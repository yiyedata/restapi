package restapi

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	default404Body = []byte("404 not found")
)

type HandlerFunc func(*Context)
type Router struct {
	funs []HandlerFunc
	k1   string
	k2   string
	k3   string
}

type Engine struct {
	prefix      string
	middlewares []HandlerFunc
	router      map[string]*Router
}

func New() *Engine {
	return &Engine{router: make(map[string]*Router)}
}
func (engine *Engine) RootGroup(prefix string, middlewares ...HandlerFunc) {
	engine.prefix = prefix
	engine.middlewares = middlewares
}

func (engine *Engine) Group(prefix string, middlewares ...HandlerFunc) *RouterGroup {
	g := &RouterGroup{prefix: prefix, middlewares: make([]HandlerFunc, 0), engine: engine}
	g.Use(middlewares...)
	return g
}
func (engine *Engine) addRoutes(method string, pattern string, handler HandlerFunc, middlewares []HandlerFunc) {
	funs := make([]HandlerFunc, 0)
	funs = append(funs, middlewares...)
	funs = append(funs, handler)
	engine.addRoute(method, pattern, funs...)
}
func (engine *Engine) addRoute(method string, pattern string, handler ...HandlerFunc) {
	pattern = engine.prefix + pattern
	var handlers []HandlerFunc
	if engine.middlewares != nil && len(engine.middlewares) > 0 {
		handlers = make([]HandlerFunc, 0)
		handlers = append(handlers, engine.middlewares...)
		handlers = append(handlers, handler...)
	} else {
		handlers = handler
	}
	start := strings.Index(pattern, "/:")
	if start >= 0 {
		key := method + "-" + pattern[0:start]
		router := &Router{}
		engine.router[key] = router
		router.funs = handlers
		paras := strings.Split(pattern[start+2:], "/:")
		router.k1 = paras[0]
		if len(paras) > 1 {
			router.k2 = paras[1]
		}
		if len(paras) > 2 {
			router.k3 = paras[2]
		}
	} else {
		key := method + "-" + pattern
		engine.router[key] = &Router{funs: handlers}
	}
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}
func (engine *Engine) DELETE(pattern string, handler HandlerFunc) {
	engine.addRoute("DELETE", pattern, handler)
}
func (engine *Engine) PUT(pattern string, handler HandlerFunc) {
	engine.addRoute("PUT", pattern, handler)
}
func (engine *Engine) OPTIONS(pattern string, handler HandlerFunc) {
	engine.addRoute("OPTIONS", pattern, handler)
}
func (engine *Engine) PATCH(pattern string, handler HandlerFunc) {
	engine.addRoute("PATCH", pattern, handler)
}
func (engine *Engine) TRACE(pattern string, handler HandlerFunc) {
	engine.addRoute("TRACE", pattern, handler)
}
func (engine *Engine) CONNECT(pattern string, handler HandlerFunc) {
	engine.addRoute("CONNECT", pattern, handler)
}
func (engine *Engine) HEAD(pattern string, handler HandlerFunc) {
	engine.addRoute("HEAD", pattern, handler)
}
func (engine *Engine) Method(method, pattern string, handler HandlerFunc) {
	engine.addRoute(method, pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	fmt.Println("API server listening at " + addr)
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	c := newContext(w, req)
	handlers := engine.getRouter(key, c)
	if handlers != nil {
		c.handlers = handlers
		c.Next()
	} else {
		c.String(404, default404Body)
	}
}
func (engine *Engine) getRouter(url string, c *Context) []HandlerFunc {
	router := engine.router
	if handler, ok := router[url]; ok {
		return handler.funs
	} else {
		k, k1 := engine.getKV(url)
		if handler, ok := router[k]; ok {
			c.SetPara(engine.router[k].k1, k1)
			return handler.funs
		} else if k != "" {
			k, k2 := engine.getKV(k)
			if k != "" {
				if handler, ok := router[k]; ok {
					c.SetPara(engine.router[k].k1, k2)
					c.SetPara(engine.router[k].k2, k1)
					return handler.funs
				} else {
					k, k3 := engine.getKV(k)
					if k != "" {
						if handler, ok := router[k]; ok {
							c.SetPara(engine.router[k].k1, k3)
							c.SetPara(engine.router[k].k2, k2)
							c.SetPara(engine.router[k].k3, k1)
							return handler.funs
						}
					}
				}
			}
		}
		return nil
	}
}
func (engine *Engine) getKV(key string) (string, string) {
	start := strings.LastIndex(key, "/")
	if start >= 0 {
		k := key[0:start]
		v := key[start+1:]
		return k, v
	}
	return "", ""
}
