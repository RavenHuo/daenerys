/**
 * @Author raven
 * @Description
 * @Date 2022/7/11
 **/
package main

import (
	"github.com/RavenHuo/daenerys/http/server"
	"github.com/RavenHuo/go-kit/log"
)

type FirstFilter struct{}

func (f *FirstFilter) Order() int {
	return 1
}
func (f *FirstFilter) Name() string {
	return "FirstFilter"
}

func (f *FirstFilter) DoFilter(c *server.Context, chain *server.HandlerFilterChain) error {
	log.Infof(c.Ctx, "hello 1")
	return chain.DoFilter(c)
}

type SecondFilter struct{}

func (s *SecondFilter) Order() int {
	return 2
}
func (s *SecondFilter) Name() string {
	return "SecondFilter"
}

func (s *SecondFilter) DoFilter(c *server.Context, chain *server.HandlerFilterChain) error {
	log.Info(c.Ctx, "hello 2")
	return chain.DoFilter(c)
}
