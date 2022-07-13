/**
 * @Author raven
 * @Description
 * @Date 2022/7/11
 **/
package intercept

import (
	"fmt"
	"math"

	"github.com/RavenHuo/daenerys/http/server"
	"github.com/RavenHuo/daenerys/internal/tls"
	"github.com/RavenHuo/daenerys/log"
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

func (l *LogIntercept) PreHandle(c *server.Context) bool {
	traceId := fmt.Sprintf(traceIdFormat, c.ServerName, tls.GoID())
	c.Ctx = context.WithValue(c.Ctx, log.TraceIdField, traceId)
	return true
}

func (l *LogIntercept) AfterCompletion(context *server.Context) {
}
