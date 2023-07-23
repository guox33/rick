package main

import (
	"github.com/guox33/rick/http_framework/framework"
	"net/http"
)

func FooControllerHandler(ctx *framework.Context) {
	_ = ctx.Json(200, map[string]interface{}{
		"hello": "foo",
	})
}

func UserLoginHandler(ctx *framework.Context) {
	_ = ctx.Json(200, map[string]interface{}{
		"hello": "user login",
	})
}

func SubjectDelHandler(ctx *framework.Context) {
	panic("del panic")
	/*_ = ctx.Json(200, map[string]interface{}{
		"hello": "subject del",
	})*/
}

func SubjectUpdateHandler(ctx *framework.Context) {
	_ = ctx.Json(200, map[string]interface{}{
		"hello": "subject update",
	})
}

func SubjectGetHandler(ctx *framework.Context) {
	_ = ctx.Json(200, map[string]interface{}{
		"hello": "subject get",
	})
}

func SubjectListHandler(ctx *framework.Context) {
	_ = ctx.Json(200, map[string]interface{}{
		"hello": "subject list",
	})
}

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := http.Server{Handler: core, Addr: ":80"}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
