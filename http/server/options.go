package server

import (
	"time"

	"github.com/opentracing/opentracing-go"
)

const (
	HTTPReadTimeout  = 60 * time.Second
	HTTPWriteTimeout = 60 * time.Second
	HTTPIdleTimeout  = 90 * time.Second
	defaultBodySize  = 1024
)

type Options struct {
	tracer             opentracing.Tracer
	port               int
	readTimeout        time.Duration
	writeTimeout       time.Duration
	idleTimeout        time.Duration // server keep conn
	certFile           string
	keyFile            string
	reqBodyLogOff      bool
	respBodyLogMaxSize int
	recoverPanic       bool
}

func DefaultOptions() *Options {
	opts := &Options{}
	if opts.readTimeout == 0 {
		opts.readTimeout = HTTPReadTimeout
	}
	if opts.writeTimeout == 0 {
		opts.writeTimeout = HTTPWriteTimeout
	}
	if opts.idleTimeout == 0 {
		opts.idleTimeout = HTTPIdleTimeout
	}
	if opts.respBodyLogMaxSize == 0 {
		opts.respBodyLogMaxSize = defaultBodySize
	}
	return opts
}

func (o *Options) Tracer(tracer opentracing.Tracer) *Options {
	if tracer != nil {
		o.tracer = tracer
	}
	return o
}

func (o *Options) Port(port int) *Options {
	o.port = port
	return o
}

// 从连接被接受(accept)到request body完全被读取(如果你不读取body，那么时间截止到读完header为止)
// 包括了TCP消耗的时间,读header时间
// 对于 https请求，ReadTimeout 包括了TLS握手的时间
func (o *Options) ReadTimeout(d time.Duration) *Options {
	o.readTimeout = d
	return o
}

// 从request header的读取结束开始，到response write结束为止 (也就是 ServeHTTP 方法的声明周期)
func (o *Options) WriteTimeout(d time.Duration) *Options {
	o.writeTimeout = d
	return o
}

func (o *Options) IdleTimeout(d time.Duration) *Options {
	o.idleTimeout = d
	return o
}

func (o *Options) CertFile(file string) *Options {
	o.certFile = file
	return o
}

func (o *Options) KeyFile(file string) *Options {
	o.keyFile = file
	return o
}

// 关闭req body 打印
func (o *Options) RequestBodyLogOff(b bool) *Options {
	o.reqBodyLogOff = b
	return o
}

// 控制resp body打印大小
func (o *Options) RespBodyLogMaxSize(size int) *Options {
	o.respBodyLogMaxSize = size
	return o
}

func (o *Options) RecoverPanic(rp bool) *Options {
	o.recoverPanic = rp
	return o
}
