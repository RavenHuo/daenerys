package server

//import (
//	"git.inke.cn/inkelogic/daenerys/http/third"
//	"github.com/RavenHuo/daenerys/core"
//)
//
//var serverInternalThirdPlugin = third.New()
//
//// plugins will be effect always
//func RegisterOnGlobalStage(plugFunc ...HandlerFunc) {
//	ps := make([]core.Plugin, len(plugFunc))
//	for i := range plugFunc {
//		ps[i] = plugFunc[i]
//	}
//	serverInternalThirdPlugin.OnGlobalStage().Register(ps)
//}
//
//// plugins will be effect for a http request or a http route
//func RegisterOnRequestStage(plugFunc ...HandlerFunc) {
//	ps := make([]HandlerIntercept, len(plugFunc))
//	for i := range plugFunc {
//		ps[i] = plugFunc[i]
//	}
//	serverInternalThirdPlugin.OnRequestStage().Register(ps)
//}
