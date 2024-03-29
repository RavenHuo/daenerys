package server

import (
	"context"
	"fmt"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/utils/nets"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RavenHuo/daenerys/internal/tls"
)

// core plugin encapsulation
type HandlerFunc func(c *RContext)

type Server interface {
	Router
	Run(addr ...string) error
	Stop() error
	Filters(p ...HandlerFilter)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type server struct {
	RouterMgr
	options      *Options
	mu           sync.Mutex
	trees        methodTrees
	srv          *http.Server
	running      int32
	once         sync.Once
	pool         sync.Pool
	paths        []string
	onHijackMode bool
	filter       []HandlerFilter
	name         string
}

func NewServer(options *Options) Server {
	context.Background()
	s := &server{
		RouterMgr: RouterMgr{
			intercepts: nil,
			basePath:   "/",
		},
		trees:  make(methodTrees, 0, 10),
		mu:     sync.Mutex{},
		pool:   sync.Pool{},
		filter: make([]HandlerFilter, 0, 2),
	}
	s.pool.New = func() interface{} {
		return s.allocContext()
	}
	s.options = options
	s.srv = &http.Server{
		Handler:      s,
		ReadTimeout:  s.options.readTimeout,
		WriteTimeout: s.options.writeTimeout,
		IdleTimeout:  s.options.idleTimeout,
		ConnState: func(conn net.Conn, state http.ConnState) {
			switch state {
			case http.StateHijacked:
				s.onHijackMode = true
			}
		},
	}
	s.RouterMgr.server = s
	s.addDefaultRoute()
	atomic.StoreInt32(&s.running, 0)

	return s
}

func (s *server) Run(addr ...string) error {
	var err error
	var host string
	s.once.Do(func() {
		s.uploadServerPath()
		port := 0
		if len(addr) > 0 {
			host = addr[0]
			tmp := strings.Split(host, ":")
			if len(tmp) == 2 {
				port, _ = strconv.Atoi(tmp[1])
			} else {
				err = fmt.Errorf("invalid addr: %s", addr)
				return
			}
		} else if s.options.port > 0 {
			port = s.options.port
			host = fmt.Sprintf(":%d", port)
		} else {
			host = ":80"
		}
		ln, e := net.Listen("tcp", host)
		if e != nil {
			fmt.Printf("start http server on %s failed, %v\n", host, e)
			err = e
			return
		}
		fmt.Printf("start http server on %s\n", host)

		if len(s.options.certFile) == 0 || len(s.options.keyFile) == 0 {
			err = s.srv.Serve(ln)
		} else {
			err = s.srv.ServeTLS(ln, s.options.certFile, s.options.keyFile)
		}
		if err != nil {
			if err == http.ErrServerClosed {
				fmt.Printf("http server closed: %v", err)
				err = nil
			}
		}
		fmt.Println("http server start success")
	})
	return err
}

func (s *server) Stop() error {
	if !atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		return nil
	}

	// gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err := s.srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		fmt.Printf("gracefully shutdown, err:%v", err)
	}
	cancel()
	return nil
}

func (s *server) allocContext() *RContext {
	return &RContext{
		srv: s,
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := s.pool.Get().(*RContext)
	ctx.reset()
	defer s.pool.Put(ctx)

	ctx.ServerName = s.name
	ctx.Ctx = context.Background()
	ctx.startTime = time.Now()
	ctx.w.reset(w, s.options.respBodyLogMaxSize)
	ctx.Request = r
	ctx.Response = ctx.w
	s.handleHTTPRequest(ctx)
}

func (s *server) handleHTTPRequest(ctx *RContext) {
	nodeValue := ctx.requestNode()
	if nodeValue.handler == nil {
		if s.methodNotAllowed(ctx) {
			ctx.Response.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = ctx.Response.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
			fmt.Printf("http server, method not allowd, request %v\n", *ctx.Request)
			return
		}

		ctx.Response.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Response.Write([]byte(http.StatusText(http.StatusNotFound)))
		fmt.Printf(" http server, handlers not found, request %+v\n", *ctx.Request)
		return
	}

	nCtx := context.Background()
	tls.SetContext(nCtx)
	defer tls.Flush()

	s.internalHandle(ctx, nodeValue)

}

// internal handler http request
func (s *server) internalHandle(ctx *RContext, nodeValue *nodeValue) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorf(ctx.Ctx, "internalHandle err :%s", err)
			ctx.Response.WriteHeader(http.StatusInternalServerError)
			_, _ = ctx.Response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}
	}()
	// 过滤器
	filterChain := MakeFilterChain(s.filter)
	if err := filterChain.DoFilter(ctx); err != nil {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		_, _ = ctx.Response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	for _, ins := range nodeValue.intercepts {
		// pre handle intercept
		if !ins.PreHandle(ctx) {
			return
		}
	}
	nodeValue.handler(ctx)
	for _, ins := range nodeValue.intercepts {
		ins.AfterCompletion(ctx)
	}
}

func (s *server) Filters(filters ...HandlerFilter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.filter = append(s.filter, filters...)
}

func (s *server) addRoute(method, path string, handler HandlerFunc, intercepts []HandlerIntercept) {
	if path[0] != '/' || len(method) == 0 || handler == nil {
		return
	}
	root := s.trees.get(method)
	if root == nil {
		root = new(node)
		s.trees = append(s.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handler, intercepts)
	s.addPath(path)
}

func (s *server) addPath(path string) {
	var exist = false
	for _, v := range s.paths {
		if v == path {
			exist = true
			break
		}
	}
	if !exist {
		s.paths = append(s.paths, path)
	}
}

func getRemoteIP(r *http.Request) string {
	for _, h := range []string{"X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := addresses[i]
			if len(ip) > 0 {
				return nets.IPFormat(ip)
			}
		}
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return nets.IPFormat(ip)
}

// TODO
func (s *server) uploadServerPath() {
	//body := map[string]interface{}{}
	//body["type"] = 1
	//body["resource_list"] = s.paths
	//body["service"] = s.options.serviceName
	//b, _ := json.NewEncoder().Encode(body)
	//respB, err := tracing.KVPut(b)
	//if err != nil {
	//	return
	//}
	//logging.GenLogf("sync http server path list to consul response:%q", respB)
}

// 判断是不是请求方式出错
func (s *server) methodNotAllowed(ctx *RContext) bool {
	// 405
	t := s.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method == ctx.Request.Method {
			continue
		}
		root := t[i].root
		// plugin, urlparam, found, matchPath expression
		nodeValue := root.getValue(ctx.Request.URL.Path, ctx.Params, false)
		// 存在 路径相同，但是请求方式不一样的node
		if nodeValue.handler != nil {
			return true
		}
	}
	return false
}

func (s *server) addDefaultRoute() {
	s.GET("/", func(c *RContext) {
		c.Response.Write([]byte("hello"))
	})
}
