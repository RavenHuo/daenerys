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
	serviceName        string
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

type Option func(*Options)

func newOptions(options ...Option) Options {
	opts := Options{}
	for _, o := range options {
		o(&opts)
	}
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

func Tracer(tracer opentracing.Tracer) Option {
	return func(o *Options) {
		if tracer != nil {
			o.tracer = tracer
		}
	}
}

func Port(port int) Option {
	return func(o *Options) {
		o.port = port
	}
}

func Name(serviceName string) Option {
	return func(o *Options) {
		o.serviceName = serviceName
	}
}

// 从连接被接受(accept)到request body完全被读取(如果你不读取body，那么时间截止到读完header为止)
// 包括了TCP消耗的时间,读header时间
// 对于 https请求，ReadTimeout 包括了TLS握手的时间
func ReadTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.readTimeout = d
	}
}

// 从request header的读取结束开始，到response write结束为止 (也就是 ServeHTTP 方法的声明周期)
func WriteTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.writeTimeout = d
	}
}

func IdleTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = d
	}
}

func CertFile(file string) Option {
	return func(o *Options) {
		o.certFile = file
	}
}

func KeyFile(file string) Option {
	return func(o *Options) {
		o.keyFile = file
	}
}

// 关闭req body 打印
func RequestBodyLogOff(b bool) Option {
	return func(o *Options) {
		o.reqBodyLogOff = b
	}
}

// 控制resp body打印大小
func RespBodyLogMaxSize(size int) Option {
	return func(o *Options) {
		o.respBodyLogMaxSize = size
	}
}

func RecoverPanic(rp bool) Option {
	return func(o *Options) {
		o.recoverPanic = rp
	}
}
