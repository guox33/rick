package main

import "github.com/guox33/rick/http_framework/framework"

func registerRouter(core *framework.Core) {
	core.Get("/foo", FooControllerHandler)
	core.Get("/user/login", UserLoginHandler)
	subjectApi := core.Group("/subject")
	{
		subjectApi.Delete("/:id", SubjectDelHandler)
		subjectApi.Put("/:id", SubjectUpdateHandler)
		subjectApi.Get("/:id", SubjectGetHandler)
		subjectApi.Get("/list/all", SubjectListHandler)
	}
}
