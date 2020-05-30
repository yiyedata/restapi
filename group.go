package restapi

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	engine      *Engine       // all groups share a Engine instance
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {

	group.engine.addRoutes("GET", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.engine.addRoutes("POST", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.engine.addRoutes("DELETE", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.engine.addRoutes("PUT", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFunc) {
	group.engine.addRoutes("OPTIONS", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.engine.addRoutes("PATCH", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) TRACE(pattern string, handler HandlerFunc) {
	group.engine.addRoutes("TRACE", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) CONNECT(pattern string, handler HandlerFunc) {
	group.engine.addRoutes("CONNECT", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	group.engine.addRoutes("HEAD", group.prefix+pattern, handler, group.middlewares)
}
func (group *RouterGroup) Method(method, pattern string, handler HandlerFunc) {
	group.engine.addRoutes(method, group.prefix+pattern, handler, group.middlewares)
}

// middlewares = append(middlewares, group.middlewares...)
