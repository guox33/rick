package middleware

import (
	"fmt"
	"github.com/guox33/rick/http_framework/framework"
)

func RecoveryMiddleware() framework.ControlHandler {
	return func(ctx *framework.Context) {
		defer func() {
			if r := recover(); r != nil {
				_ = ctx.Text(500, fmt.Sprintf("panic: %v", r))
			}
		}()

		ctx.Next()
	}
}
