package main

import (
	"code.byted.org/clientQA/rick/http_framework/framework"
	"net/http"
)

func FooHandle(request *http.Request, response http.ResponseWriter) {
	ctx := framework.NewContext(request, response)
	_ = FooControllerHandler(ctx)
}

func FooControllerHandler(ctx *framework.Context) error {
	return ctx.Json(200, map[string]interface{}{
		"hello": "world",
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
