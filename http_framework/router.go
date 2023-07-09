package main

import "code.byted.org/clientQA/rick/http_framework/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooControllerHandler)
}
