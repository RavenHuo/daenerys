/**
 * @Author raven
 * @Description
 * @Date 2022/6/26
 **/
package core

import (
	"context"
)

type Core interface {
	Use(...Intercept) Core
	Next(context.Context)
	AbortErr(error)
	Abort()
	IsAborted() bool
	Err() error
	Copy() Core
	Index() int
	Reset(idx int)
}