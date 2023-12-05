/**
 * @Author raven
 * @Description
 * @Date 2022/7/11
 **/
package intercept

import (
	"fmt"
	"github.com/RavenHuo/go-pkg/trace"
	"math"

	"github.com/RavenHuo/daenerys/http/server"
	"github.com/RavenHuo/daenerys/internal/tls"
	"golang.org/x/net/context"
)

const traceIdFormat = "s_%s:t_%d"

type LogIntercept struct {
}

func (l *LogIntercept) Order() int {
	return math.MaxInt32
}

func (l *LogIntercept) Name() string {
	return "LogIntercept"
}

func (l *LogIntercept) PreHandle(c *server.RContext) bool {
	traceId := fmt.Sprintf(traceIdFormat, c.ServerName, tls.GoID())
	c.Ctx = context.WithValue(c.Ctx, trace.TraceIdField, traceId)
	return true
}

func (l *LogIntercept) AfterCompletion(context *server.RContext) {
}
