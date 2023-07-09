package framework

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Core represent core struct
type Core struct {
	router map[string]ControllerHandler

	requestTimeout time.Duration
}

func NewCore() *Core {
	return &Core{
		router:         map[string]ControllerHandler{},
		requestTimeout: 2 * time.Second,
	}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.serveHTTP")
	ctx := NewContext(request, response)

	durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), c.requestTimeout)
	defer func() {
		cancel()
	}()

	// 一个简单的路由选择器，这里直接写死为测试路由foo
	router := c.router["foo"]
	if router == nil {
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
