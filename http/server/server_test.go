package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func InitTestServer(handlers map[string]HandlerFunc, serverPort int) {
	// init tracer
	//cfg := jaegerconfig.Configuration{
	//	// SamplingServerURL: "http://localhost:5778/sampling"
	//	Sampler: &jaegerconfig.SamplerConfig{Type: jaeger.SamplerTypeRemote},
	//	Reporter: &jaegerconfig.ReporterConfig{
	//		LogSpans:            false,
	//		BufferFlushInterval: 1 * time.Second,
	//		LocalAgentHostPort:  "127.0.0.1:6831",
	//	},
	//}
	//tracer, _, err := cfg.New("danerys.test.service")
	//if err != nil {
	//	panic(err)
	//}

	// server
	s := NewServer(Name("danerys.test.service"))

	for k, v := range handlers {
		s.GET(k, v)
	}

	go func() {
		err := s.Run(fmt.Sprintf(":%d", serverPort))
		if err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(1 * time.Second)
}

func TestHttpServer_WriteHeader(t *testing.T) {
	getJson := BaseResp{
		Code: 0,
		Msg:  "操作成功",
		Body: struct{ Action string }{Action: "get json"},
	}

	InitTestServer(
		map[string]HandlerFunc{
			"/json/get/500": func(c *Context) {
				c.Response.WriteHeader(500)
				c.JSON(getJson)
				return
			},
			"/json/get/502": func(c *Context) {
				c.Response.WriteHeader(502)
				c.JSON(map[string]interface{}{"action": "get json"})
				return
			},
			"/json/get/400": func(c *Context) {
				c.Response.WriteHeader(400)
				b, _ := json.Marshal(getJson)
				_, _ = c.Response.Write(b)
				return
			},
			"/json/post/403": func(c *Context) {
				c.Response.WriteHeader(403)
				_, _ = c.Response.WriteString("hello world 403")
				return
			},
			"/json/post/200": func(c *Context) {
				c.Response.WriteHeader(200)
				c.JSON(map[string]interface{}{"action": "post json"})
				return
			},
			"/add/header": func(c *Context) {
				c.Response.Header().Add("x-my-header-1", "hello world 1")
				c.Response.Header().Add("x-my-header-2", "hello world 2")
				c.Response.Header().Add("x-my-header-3", "hello world 3")
				c.Response.WriteHeader(200)
				_, _ = c.Response.WriteString("add hello world header")
				return
			},
		},
		22356,
	)

	httpclient := http.Client{Timeout: 10 * time.Second}

	fmt.Println("========= 500 header =========")
	// 500 status
	rsp, err := httpclient.Get("http://localhost:22356/json/get/500")
	assert.Equal(t, nil, err)
	if rsp == nil {
		t.Fail()
	}
	respB, err := ioutil.ReadAll(rsp.Body)
	assert.Equal(t, nil, err)
	jsonGet500Response := BaseResp{}
	err = json.Unmarshal(respB, &jsonGet500Response)
	assert.Equal(t, nil, err)
	assert.Equal(t, getJson, jsonGet500Response)

	assert.Equal(t, "500 Internal Server Error", rsp.Status)

	fmt.Println("========= 502 header =========")
	// 502 status
	rsp, err = httpclient.Get("http://localhost:22356/json/get/502")
	assert.Equal(t, nil, err)
	if rsp == nil {
		t.Fail()
	}
	respB, err = ioutil.ReadAll(rsp.Body)
	assert.Equal(t, nil, err)
	jsonGet502Response := BaseResp{}
	err = json.Unmarshal(respB, &jsonGet502Response)
	assert.Equal(t, nil, err)
	assert.Equal(t, BaseResp{Code: 0, Msg: "0", Body: struct{ Action string }{Action: "get json"}}, jsonGet502Response)

	assert.Equal(t, "502 Bad Gateway", rsp.Status)

	fmt.Println("========= 400 header =========")
	// 400 status
	rsp, err = httpclient.Get("http://localhost:22356/json/get/400")
	assert.Equal(t, nil, err)
	if rsp == nil {
		t.Fail()
	}
	respB, err = ioutil.ReadAll(rsp.Body)
	assert.Equal(t, nil, err)
	jsonGet400Response := BaseResp{}
	err = json.Unmarshal(respB, &jsonGet400Response)
	assert.Equal(t, nil, err)
	assert.Equal(t, getJson, jsonGet400Response)

	assert.Equal(t, "400 Bad Request", rsp.Status)

	fmt.Println("========= 403 header =========")
	// 403 status
	rsp, err = httpclient.Post("http://localhost:22356/json/post/403", "Content-Type: application/json; charset=utf-8", nil)
	assert.Equal(t, nil, err)
	if rsp == nil {
		t.Fail()
	}
	respB, err = ioutil.ReadAll(rsp.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "hello world 403", string(respB[:]))

	assert.Equal(t, "403 Forbidden", rsp.Status)

	fmt.Println("========= 200 header =========")
	// 200 status
	rsp, err = httpclient.Post("http://localhost:22356/json/post/200", "Content-Type: application/json; charset=utf-8", nil)
	assert.Equal(t, nil, err)
	if rsp == nil {
		t.Fail()
	}
	respB, err = ioutil.ReadAll(rsp.Body)
	assert.Equal(t, nil, err)
	jsonPost200Response := BaseResp{}
	err = json.Unmarshal(respB, &jsonPost200Response)
	assert.Equal(t, nil, err)
	assert.Equal(t, BaseResp{Code: 0, Msg: "0", Body: struct{ Action string }{Action: "post json"}}, jsonPost200Response)

	fmt.Println("========= add header =========")
	// add header
	rsp, err = httpclient.Post("http://localhost:22356/add/header", "Content-Type: application/json; charset=utf-8", nil)
	assert.Equal(t, nil, err)
	if rsp == nil {
		t.Fail()
	}
	respB, err = ioutil.ReadAll(rsp.Body)
	assert.Equal(t, nil, err)
	assert.Equal(t, "add hello world header", string(respB[:]))

	assert.Equal(t, "200 OK", rsp.Status)
	assert.Equal(t, "hello world 1", rsp.Header.Get("x-my-header-1"))
	assert.Equal(t, "hello world 2", rsp.Header.Get("x-my-header-2"))
	assert.Equal(t, "hello world 3", rsp.Header.Get("x-my-header-3"))
	if rsp.Header.Get("X-Trace-Id") == "" {
		t.Fail()
	}
}

func Test_getRemoteIP(t *testing.T) {
	newRequest := func(addr string) (req *http.Request) {
		req = &http.Request{
			Header: http.Header{},
		}
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			req.Header.Add("X-Real-Ip", addr)
			return
		}
		if len(port) > 0 {
			req.RemoteAddr = addr
			req.Header.Add("X-Real-Ip", host)
			return
		}
		req.Header.Add("X-Real-Ip", addr)
		return
	}

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "192.168.32.48",
			args: args{
				r: newRequest("192.168.32.48"),
			},
			want: "192.168.32.48",
		},
		{
			name: "192.168.32.48:7085",
			args: args{
				r: newRequest("192.168.32.48:7085"),
			},
			want: "192.168.32.48",
		},
		{
			name: "ali-c-inf-testing01.bj",
			args: args{
				r: newRequest("ali-c-inf-testing01.bj"),
			},
			want: "ali-c-inf-testing01.bj",
		},
		{
			name: "ali-c-inf-testing01.bj:7085",
			args: args{
				r: newRequest("ali-c-inf-testing01.bj:7085"),
			},
			want: "ali-c-inf-testing01.bj",
		},

		{
			name: "1050:0000:0000:0000:0005:0600:300c:326b",
			args: args{
				r: newRequest("1050:0000:0000:0000:0005:0600:300c:326b"),
			},
			want: "1050:0:0:0:5:600:300c:326b",
		},
		{
			name: "1050:0:0:0:5:600:300c:326b",
			args: args{
				r: newRequest("1050:0:0:0:5:600:300c:326b"),
			},
			want: "1050:0:0:0:5:600:300c:326b",
		},
		{
			name: "ff06:0:0:0:0:0:0:c3",
			args: args{
				r: newRequest("ff06:0:0:0:0:0:0:c3"),
			},
			want: "ff06:0:0:0:0:0:0:c3",
		},
		{
			name: "ff06::c3",
			args: args{
				r: newRequest("ff06::c3"),
			},
			want: "ff06:0:0:0:0:0:0:c3",
		},
		{
			name: "::00c3",
			args: args{
				r: newRequest("::00c3"),
			},
			want: "0:0:0:0:0:0:0:c3",
		},
		{
			name: "00c3::",
			args: args{
				r: newRequest("00c3::"),
			},
			want: "c3:0:0:0:0:0:0:0",
		},
		{
			name: "0:0:0:0:0:ffff:192.1.56.10",
			args: args{
				r: newRequest("0:0:0:0:0:ffff:192.1.56.10"),
			},
			want: "0:0:0:0:0:ffff:192.1.56.10",
		},
		{
			name: "::192.1.56.10",
			args: args{
				r: newRequest("::192.1.56.10"),
			},
			want: "0:0:0:0:0:0:192.1.56.10",
		},
		{
			name: "::ffff:192.1.56.10",
			args: args{
				r: newRequest("::ffff:192.1.56.10"),
			},
			want: "0:0:0:0:0:ffff:192.1.56.10",
		},
		{
			name: "3b00:2cb1::192.1.56.10",
			args: args{
				r: newRequest("3b00:2cb1::192.1.56.10"),
			},
			want: "3b00:2cb1:0:0:0:0:192.1.56.10",
		},
		{
			name: "2408:400a:16d:6e00:d18:a9d:5feb:e6f9",
			args: args{
				r: newRequest("2408:400a:16d:6e00:d18:a9d:5feb:e6f9"),
			},
			want: "2408:400a:16d:6e00:d18:a9d:5feb:e6f9",
		},
		{
			name: "[ff06::c3]:7085",
			args: args{
				r: newRequest("[ff06::c3]:7085"),
			},
			want: "ff06:0:0:0:0:0:0:c3",
		},
		{
			name: "[::00c3]:7085",
			args: args{
				r: newRequest("[::00c3]:7085"),
			},
			want: "0:0:0:0:0:0:0:c3",
		},
		{
			name: "[00c3::]:7085",
			args: args{
				r: newRequest("[00c3::]:7085"),
			},
			want: "c3:0:0:0:0:0:0:0",
		},
		{
			name: "[::ffff:192.1.56.10]:7085",
			args: args{
				r: newRequest("[::ffff:192.1.56.10]:7085"),
			},
			want: "0:0:0:0:0:ffff:192.1.56.10",
		},
		{
			name: "[2408:400a:16d:6e00:d18:a9d:5feb:e6f9]:7085",
			args: args{
				r: newRequest("[2408:400a:16d:6e00:d18:a9d:5feb:e6f9]:7085"),
			},
			want: "2408:400a:16d:6e00:d18:a9d:5feb:e6f9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRemoteIP(tt.args.r); got != tt.want {
				t.Errorf("getRemoteIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
