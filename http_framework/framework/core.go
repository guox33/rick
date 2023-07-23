package framework

import (
	"log"
	"net/http"
	"strings"
)

// Core represent core struct
type Core struct {
	router map[string]*Tree

	middlewares ControlHandlerChain
}

func NewCore() *Core {
	return &Core{
		router: map[string]*Tree{
			http.MethodGet:    NewTree(),
			http.MethodPost:   NewTree(),
			http.MethodPut:    NewTree(),
			http.MethodDelete: NewTree(),
		},
	}
}

func (c *Core) Get(url string, handler ...ControlHandler) {
	allHandler := append(c.middlewares, handler...)
	if err := c.router[http.MethodGet].AddRouter(strings.ToUpper(url), allHandler...); err != nil {
		log.Fatalf("add router error: %v", err)
	}
}

func (c *Core) Post(url string, handler ...ControlHandler) {
	allHandler := append(c.middlewares, handler...)
	if err := c.router[http.MethodPost].AddRouter(strings.ToUpper(url), allHandler...); err != nil {
		log.Fatalf("add router error: %v", err)
	}
}

func (c *Core) Delete(url string, handler ...ControlHandler) {
	allHandler := append(c.middlewares, handler...)
	if err := c.router[http.MethodDelete].AddRouter(strings.ToUpper(url), allHandler...); err != nil {
		log.Fatalf("add router error: %v", err)
	}
}

func (c *Core) Put(url string, handler ...ControlHandler) {
	allHandler := append(c.middlewares, handler...)
	if err := c.router[http.MethodPut].AddRouter(strings.ToUpper(url), allHandler...); err != nil {
		log.Fatalf("add router error: %v", err)
	}
}

func (c *Core) Group(url string, handler ...ControlHandler) IGroup {
	return NewGroup(c, nil, url, append(c.middlewares, handler...)...)
}

func (c *Core) Use(middlewares ...ControlHandler) {
	c.middlewares = middlewares
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	ctx := NewContext(request, response)

	// 一个简单的路由选择器，这里直接写死为测试路由foo
	router := c.findNodeByRequest(request)
	if router == nil {
		_ = ctx.Text(404, "not found")
		return
	}
	log.Println("core.router")

	ctx.handlerChain = router.handler
	ctx.Next()
}

func (c *Core) findNodeByRequest(request *http.Request) *node {
	if router, ok := c.router[strings.ToUpper(request.Method)]; ok {
		return router.findNode(request.URL.Path)
	}
	return nil
}
