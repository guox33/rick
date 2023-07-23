package framework

type ControlHandler func(c *Context)

type ControlHandlerChain []ControlHandler

type Middleware func(next ControlHandler) ControlHandler
