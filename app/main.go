/**
 * @Author raven
 * @Description
 * @Date 2022/7/11
 **/
package main

import (
	"github.com/RavenHuo/daenerys/app/intercept"
	"github.com/RavenHuo/daenerys/http/server"
)

func main() {
	firstFilter := &FirstFilter{}
	secondFilter := &SecondFilter{}
	options := server.DefaultOptions().Name("test")
	httpServer := server.NewServer(options)
	httpServer.GET("/ping", ping, &intercept.FirstHandlerIntercept{}, &intercept.LogIntercept{})
	httpServer.Filters(firstFilter, secondFilter)
	httpServer.Run("0.0.0.0:8080")
}

func ping(c *server.Context) {
	c.JSON(server.BaseResp{
		Body: "heartbeat",
	})
}
