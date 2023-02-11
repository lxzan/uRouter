package uRouter

import (
	"github.com/lxzan/uRouter/internal"
	"strings"
)

// Group 路由组
// route group
type Group struct {
	router      *Router
	path        string
	middlewares []HandlerFunc
}

// Group 创建子路由组
// create a child route group
func (c *Group) Group(path string, middlewares ...HandlerFunc) *Group {
	c.router.mu.Lock()
	defer c.router.mu.Unlock()

	group := &Group{
		router:      c.router,
		path:        internal.JoinPath(SEP, c.path, path),
		middlewares: append(c.middlewares, middlewares...),
	}
	return group
}

// OnEvent 监听事件
// listen to event
func (c *Group) OnEvent(action string, path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	router := c.router
	router.mu.Lock()
	defer router.mu.Unlock()

	action = strings.ToLower(action)
	h := append(c.middlewares, middlewares...)
	h = append(h, handler)
	router.apis = append(router.apis, &apiHandler{
		Action:   action,
		Path:     internal.JoinPath(SEP, c.path, path),
		FullPath: internal.JoinPath(SEP, action, c.path, path),
		Funcs:    h,
	})

	//p := internal.JoinPath(c.separator, action, c.path, path)
	//h := append(c.middlewares, middlewares...)
	//h = append(h, handler)
	//
	//// 检测路径冲突
	//if router.pathExists(p) {
	//	router.showPathConflict(p)
	//	return
	//}
	//
	//if !hasVar(p) {
	//	router.staticRoutes[p] = h
	//} else {
	//	router.dynamicRoutes.Set(p, h)
	//}
}

// On  类似OnEvent方法, 但是没有动作修饰词
func (c *Group) On(path string, handler HandlerFunc, middlewares ...HandlerFunc) {
	c.OnEvent("", path, handler, middlewares...)
}
