/**
 * @Author raven
 * @Description
 * @Date 2022/6/26
 **/
package core

import (
	"context"
)

// http/grpc 拦截器
type Intercept interface {
	Do(context.Context, Core)
}

type Function func(context.Context, Core)

func (f Function) Do(ctx context.Context, c Core) { f(ctx, c) }
