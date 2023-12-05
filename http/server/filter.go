/**
 * @Author raven
 * @Description
 * @Date 2022/7/7
 **/
package server

import (
	"errors"
	"fmt"
	"sort"

	"github.com/RavenHuo/daenerys/core"
)

// 过滤器
type HandlerFilter interface {
	core.Order
	DoFilter(*RContext, *HandlerFilterChain) error
}

func sortHandlerFilter(filters []HandlerFilter) {
	sort.Slice(filters, func(i, j int) bool {
		return filters[i].Order() > filters[j].Order()
	})
}

func MakeFilterChain(filters []HandlerFilter) *HandlerFilterChain {
	sortHandlerFilter(filters)
	return &HandlerFilterChain{
		pos:     -1,
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

func (chain *HandlerFilterChain) DoFilter(ctx *RContext) error {
	if len(chain.filters) == 0 || chain.pos == len(chain.filters)-1 {
		return nil
	}
	chain.pos += 1
	err := chain.internalDoFilter(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (chain *HandlerFilterChain) internalDoFilter(ctx *RContext) error {
	filter := chain.filters[chain.pos]
	err := filter.DoFilter(ctx, chain)
	if err != nil {
		return errors.New(fmt.Sprintf("internalDoFilter:%s filter err %s", filter.Name(), err))
	}
	return nil
}
