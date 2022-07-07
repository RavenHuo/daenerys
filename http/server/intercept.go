/**
 * @Author raven
 * @Description
 * @Date 2022/6/26
 **/
package server

import (
	"errors"
	"fmt"
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

// 过滤器
type HandlerFilter interface {
	core.Order
	DoFilter(*Context, *HandlerFilterChain) error
}

func sortHandlerFilter(filters []HandlerFilter) {
	sort.Slice(filters, func(i, j int) bool {
		return filters[i].Order() > filters[j].Order()
	})
}

func MakeFilterChain(filters []HandlerFilter) *HandlerFilterChain {
	sortHandlerFilter(filters)
	return &HandlerFilterChain{
		pos:     0,
		filters: filters,
	}
}

type HandlerFilterChain struct {
	/**
	 * The int which is used to maintain the current position
	 * in the filter requestNode.
	 */
	pos     int
	filters []HandlerFilter
}

func (chain *HandlerFilterChain) DoFilter(ctx *Context) error {
	if len(chain.filters) == 0 || chain.pos == len(chain.filters)-1 {
		return nil
	}
	return chain.internalDoFilter(ctx)
}

func (chain *HandlerFilterChain) internalDoFilter(ctx *Context) error {
	filter := chain.filters[chain.pos]
	err := filter.DoFilter(ctx, chain)
	if err != nil {
		return errors.New(fmt.Sprintf("internalDoFilter:%s filter err %s", filter.Name(), err))
	}
	chain.pos += 1
	return nil
}
