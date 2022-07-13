package server

import (
	"path"
	"sync"
)

// router
type Router interface {
	GROUP(string, ...HandlerIntercept) *RouterMgr
	GET(string, HandlerFunc, ...HandlerIntercept) Router
	POST(string, HandlerFunc, ...HandlerIntercept) Router
	DELETE(string, HandlerFunc, ...HandlerIntercept) Router
	PATCH(string, HandlerFunc, ...HandlerIntercept) Router
	PUT(string, HandlerFunc, ...HandlerIntercept) Router
	OPTIONS(string, HandlerFunc, ...HandlerIntercept) Router
	HEAD(string, HandlerFunc, ...HandlerIntercept) Router
	// TODO use 拦截器
}

type RouterMgr struct {
	intercepts  []HandlerIntercept
	handlerFunc HandlerFunc
	basePath    string
	server      *server
	mu          sync.Mutex
}

// 实现Router
func (mgr *RouterMgr) GROUP(relativePath string, intercepts ...HandlerIntercept) *RouterMgr {
	ps := make([]HandlerIntercept, len(intercepts))
	for i, h := range intercepts {
		ps[i] = h
	}
	return &RouterMgr{
		intercepts: mgr.combineHandlerIntercepts(ps),
		basePath:   mgr.absolutePath(relativePath),
		server:     mgr.server,
		mu:         sync.Mutex{},
	}
}

func (mgr *RouterMgr) POST(relativePath string, handler HandlerFunc, intercepts ...HandlerIntercept) Router {
	return mgr.handle("POST", relativePath, handler, intercepts...)
}

func (mgr *RouterMgr) GET(relativePath string, handler HandlerFunc, intercepts ...HandlerIntercept) Router {
	return mgr.handle("GET", relativePath, handler, intercepts...)
}

func (mgr *RouterMgr) DELETE(relativePath string, handler HandlerFunc, intercepts ...HandlerIntercept) Router {
	return mgr.handle("DELETE", relativePath, handler, intercepts...)
}

func (mgr *RouterMgr) PATCH(relativePath string, handler HandlerFunc, intercepts ...HandlerIntercept) Router {
	return mgr.handle("PATCH", relativePath, handler, intercepts...)
}

func (mgr *RouterMgr) PUT(relativePath string, handler HandlerFunc, intercepts ...HandlerIntercept) Router {
	return mgr.handle("PUT", relativePath, handler, intercepts...)
}

func (mgr *RouterMgr) OPTIONS(relativePath string, handler HandlerFunc, intercepts ...HandlerIntercept) Router {
	return mgr.handle("OPTIONS", relativePath, handler, intercepts...)
}

func (mgr *RouterMgr) HEAD(relativePath string, handler HandlerFunc, intercepts ...HandlerIntercept) Router {
	return mgr.handle("HEAD", relativePath, handler, intercepts...)
}

// router internal func
func (mgr *RouterMgr) handle(httpMethod, relativePath string, handler HandlerFunc, intercepts ...HandlerIntercept) Router {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	absolutePath := mgr.absolutePath(relativePath)
	his := make([]HandlerIntercept, len(intercepts))
	for i, h := range intercepts {
		his[i] = h
	}
	handlerIntercepts := mgr.combineHandlerIntercepts(his)
	sortHandlerIntercept(handlerIntercepts)
	// server has a tree, to add plugins
	mgr.server.addRoute(httpMethod, absolutePath, handler, handlerIntercepts)
	return mgr
}

// 合并拦截器
// group handlers  + handler
func (mgr *RouterMgr) combineHandlerIntercepts(handlers []HandlerIntercept) []HandlerIntercept {
	finalSize := len(mgr.intercepts) + len(handlers)
	mergedHandlers := make([]HandlerIntercept, finalSize)
	copy(mergedHandlers, mgr.intercepts)
	copy(mergedHandlers[len(mgr.intercepts):], handlers)
	return mergedHandlers
}

func (mgr *RouterMgr) absolutePath(relativePath string) string {
	if relativePath == "" || relativePath == mgr.basePath {
		return mgr.basePath
	}
	// path.join总是会把末尾/去掉
	finalPath := path.Join(mgr.basePath, relativePath)
	if relativePath[len(relativePath)-1] == '/' {
		return finalPath + "/"
	}
	return finalPath
}
