package framework

type IGroup interface {
	Get(string, ...ControlHandler)
	Post(string, ...ControlHandler)
	Put(string, ...ControlHandler)
	Delete(string, ...ControlHandler)
	Group(string, ...ControlHandler) IGroup
}

type Group struct {
	prefix      string
	parent      *Group
	core        *Core
	middlewares ControlHandlerChain
}

func NewGroup(core *Core, parent *Group, prefix string, middleware ...ControlHandler) *Group {
	return &Group{
		core:        core,
		parent:      parent,
		prefix:      prefix,
		middlewares: middleware,
	}
}

func (g *Group) Get(path string, handler ...ControlHandler) {
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Get(g.prefix+path, allHandlers...)
}

func (g *Group) Post(path string, handler ...ControlHandler) {
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Post(g.prefix+path, allHandlers...)
}

func (g *Group) Put(path string, handler ...ControlHandler) {
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Put(g.prefix+path, allHandlers...)
}

func (g *Group) Delete(path string, handler ...ControlHandler) {
	allHandlers := append(g.getMiddlewares(), handler...)
	g.core.Delete(g.prefix+path, allHandlers...)
}

func (g *Group) Group(prefix string, middleware ...ControlHandler) IGroup {
	return NewGroup(g.core, g, g.prefix+prefix, middleware...)
}

func (g *Group) Use(middlewares ...ControlHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) getMiddlewares() ControlHandlerChain {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}
