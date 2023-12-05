/**
 * @Author raven
 * @Description
 * @Date 2022/7/11
 **/
package intercept

import (
	"github.com/RavenHuo/go-pkg/log"
	"net/http"

	"github.com/RavenHuo/daenerys/http/server"
)

type FirstHandlerIntercept struct {
}

func (f *FirstHandlerIntercept) Order() int {
	return 1
}

func (f *FirstHandlerIntercept) Name() string {
	return "FirstHandlerIntercept"
}

func (f FirstHandlerIntercept) PreHandle(context *server.RContext) bool {
	context.Response.WriteHeader(http.StatusUnauthorized)
	_, _ = context.Response.Write([]byte(http.StatusText(http.StatusUnauthorized)))
	log.Info(context.Ctx, "401")
	return false
}

func (f FirstHandlerIntercept) AfterCompletion(context *server.RContext) {
	log.Info(context.Ctx, "FirstHandlerIntercept AfterCompletion")
}
