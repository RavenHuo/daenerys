package server

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/RavenHuo/daenerys/http/binding"
	"golang.org/x/net/context"
)

type RContext struct {
	Request      *http.Request
	Response     Responser
	Params       Params
	Path         string             // raw match path
	Ctx          context.Context    // for trace or others store
	intercepts   []HandlerIntercept // 拦截器
	HandlerFunc  HandlerFunc        // 处理方法
	w            *responseWriter
	loggingExtra map[string]interface{}
	// use url.ParseQuery cached the param query result from c.Request.URL.Query()
	queryCache    url.Values
	bodyBuff      *bytes.Buffer
	printRespBody bool
	printReqBody  bool
	srv           *server
	ServerName    string
	startTime     time.Time
	Keys          map[string]interface{}
	simpleBaggage map[string]string
}

func (c *RContext) reset() {
	c.Request = nil
	c.Response = nil
	c.Params = c.Params[0:0]
	c.Path = ""
	c.w = &responseWriter{}
	c.loggingExtra = nil
	c.queryCache = nil
	c.bodyBuff = bytes.NewBuffer(nil)
	c.printRespBody = true
	c.printReqBody = true
	c.Keys = nil
	c.simpleBaggage = nil
}

func (c *RContext) requestNode() *nodeValue {
	t := c.srv.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != c.Request.Method {
			continue
		}
		root := t[i].root
		// plugin, urlparam, found, matchPath expression
		v := root.getValue(c.Request.URL.Path, c.Params, false)
		return v
	}
	return nil
}

func (c *RContext) writeHeaderOnce() {
	c.Response.writeHeaderOnce()
}

func (c *RContext) LoggingExtra(k string, v interface{}) {
	if c.loggingExtra == nil {
		c.loggingExtra = map[string]interface{}{}
	}

	c.loggingExtra[k] = v

}

func (c *RContext) Bind(model interface{}) error {
	return binding.Default(c.Request, model)
}

func (c *RContext) BindJson(model interface{}) error {
	return binding.WithType(c.Request, model, binding.BindJson)
}

func (c *RContext) BindUri(model interface{}) error {
	return binding.WithType(c.Request, model, binding.BindUri)
}

// write response, response code 200
func (c *RContext) JSON(data interface{}) {
	c.Response.WriteHeader(http.StatusOK)
	_, _ = c.Response.WriteJSON(data)
}

func (c *RContext) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

// Query returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /path?id=1234&name=Manu&value=
// 	   c.Query("id") == "1234"
// 	   c.Query("name") == "Manu"
// 	   c.Query("value") == ""
// 	   c.Query("wtf") == ""
func (c *RContext) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

func (c *RContext) QueryInt(key string) int {
	i, _ := strconv.Atoi(c.Query(key))
	return i
}

func (c *RContext) QueryInt64(key string) int64 {
	i, _ := strconv.ParseInt(c.Query(key), 10, 64)
	return i
}

func (c *RContext) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *RContext) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

func (c *RContext) GetQueryArray(key string) ([]string, bool) {
	if c.queryCache == nil {
		c.queryCache, _ = url.ParseQuery(c.Request.URL.RawQuery)
	}

	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		return values, true
	}
	return []string{}, false
}

func (c *RContext) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

func (c *RContext) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

func (c *RContext) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

func (c *RContext) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *RContext) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *RContext) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *RContext) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *RContext) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *RContext) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *RContext) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *RContext) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *RContext) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

func (c *RContext) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *RContext) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
