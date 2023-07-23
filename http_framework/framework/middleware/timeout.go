package middleware

import (
	"context"
	"fmt"
	"github.com/guox33/rick/http_framework/framework"
	"log"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) framework.ControlHandler {
	return func(ctx *framework.Context) {
		durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), timeout)
		defer func() {
			cancel()

		}()

		finishCh := make(chan struct{}, 1)
		panicCh := make(chan string, 1)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					panicCh <- fmt.Sprintf("%v", r)
				}
			}()

			ctx.Next()
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
			panic(panicMsg)
		}

		return
	}
}
