package framework

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// Core represent core struct
type Core struct {
	router map[string]*Tree

	requestTimeout time.Duration
}

func NewCore() *Core {
	return &Core{
		router:         map[string]*Tree{},
		requestTimeout: 2 * time.Second,
	}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	if _, ok := c.router[http.MethodGet]; !ok {
		c.router[http.MethodGet] = NewTree()
	}
	if err := c.router[http.MethodGet].AddRouter(url, handler); err != nil {
		log.Fatalf("add router error: %v", err)
	}
}

func (c *Core) Post(url string, handler ControllerHandler) {
	if _, ok := c.router[http.MethodPost]; !ok {
		c.router[http.MethodPost] = NewTree()
	}
	if err := c.router[http.MethodPost].AddRouter(url, handler); err != nil {
		log.Fatalf("add router error: %v", err)
	}
}

func (c *Core) Delete(url string, handler ControllerHandler) {
	if _, ok := c.router[http.MethodDelete]; !ok {
		c.router[http.MethodDelete] = NewTree()
	}
	if err := c.router[http.MethodDelete].AddRouter(url, handler); err != nil {
		log.Fatalf("add router error: %v", err)
	}
}

func (c *Core) Put(url string, handler ControllerHandler) {
	if _, ok := c.router[http.MethodPut]; !ok {
		c.router[http.MethodPut] = NewTree()
	}
	if err := c.router[http.MethodPut].AddRouter(url, handler); err != nil {
		log.Fatalf("add router error: %v", err)
	}
}

func (c *Core) Group(url string) IGroup {
	return NewGroup(c, url)
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	ctx := NewContext(request, response)

	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), c.requestTimeout)
	defer func() {
		cancel()
	}()

	// 一个简单的路由选择器，这里直接写死为测试路由foo
	router := c.findRouteByRequest(request)
	if router == nil {
		_ = ctx.Text(404, "not found")
		return
	}
	log.Println("core.router")

	finishCh := make(chan struct{}, 1)
	panicCh := make(chan string, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicCh <- fmt.Sprintf("%v", r)
			}
		}()

		_ = router(ctx)
		// log.Println("done")
		finishCh <- struct{}{}
	}()

	select {
	case <-durationCtx.Done():
		log.Println("timeout")
		ctx.WriteMux().Lock()
		defer ctx.WriteMux().Unlock()
		_ = ctx.Text(500, "timeout")
		ctx.SetHasTimeout()
	case <-finishCh:

		log.Println("done")
	case panicMsg := <-panicCh:
		log.Println("panic: ", panicMsg)
		ctx.WriteMux().Lock()
		defer ctx.WriteMux().Unlock()
		_ = ctx.Text(500, "panic")
	}
}

func (c *Core) findRouteByRequest(request *http.Request) ControllerHandler {
	if router, ok := c.router[strings.ToUpper(request.Method)]; ok {
		return router.FindHandler(request.URL.Path)
	}
	return nil
}
