package framework

type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)
	Group(string) IGroup
}

type Group struct {
	prefix string
	core   *Core
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

func (g *Group) Get(path string, handler ControllerHandler) {
	g.core.Get(g.prefix+path, handler)
}

func (g *Group) Post(path string, handler ControllerHandler) {
	g.core.Post(g.prefix+path, handler)
}

func (g *Group) Put(path string, handler ControllerHandler) {
	g.core.Put(g.prefix+path, handler)
}

func (g *Group) Delete(path string, handler ControllerHandler) {
	g.core.Delete(g.prefix+path, handler)
}

func (g *Group) Group(prefix string) IGroup {
	return NewGroup(g.core, g.prefix+prefix)
}
