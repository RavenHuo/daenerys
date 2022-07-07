/**
 * @Author raven
 * @Description
 * @Date 2022/6/26
 **/
package server

import (
	"fmt"

	"github.com/RavenHuo/daenerys/errors"
)

func NewWrapResp(data interface{}, err error) BaseResp {
	e := errors.Cause(err)
	return BaseResp{
		Code: e.Code(),
		Msg:  e.Message(),
		Body: data,
	}
}

type BaseResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Body interface{} `json:"body"`
}

type query func(string) string

func exec(name string, vs ...query) string {
	ch := make(chan string)
	fn := func(i int) {
		ch <- vs[i](name)
	}
	for i, _ := range vs {
		go fn(i)
	}
	return <-ch
}

func main() {
	ret := exec("111", func(n string) string {
		return n + "func1"
	}, func(n string) string {
		return n + "func2"
	}, func(n string) string {
		return n + "func3"
	}, func(n string) string {
		return n + "func4"
	})
	fmt.Println(ret)
	// 下面还有很多执行流程，省略
}
