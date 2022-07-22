/**
 * @Author raven
 * @Description
 * @Date 2022/6/26
 **/
package server

import (
	"sort"

	"github.com/RavenHuo/daenerys/core"
)

// http/grpc 拦截器
type HandlerIntercept interface {
	core.Order
	// 前置处理
	PreHandle(*Context) bool
	// 后置处理
	AfterCompletion(*Context)
}

func sortHandlerIntercept(intercepts []HandlerIntercept) {
	sort.Slice(intercepts, func(i, j int) bool {
		return intercepts[i].Order() > intercepts[j].Order()
	})
}
