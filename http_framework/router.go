package main

import (
	"github.com/guox33/rick/http_framework/framework"
	"github.com/guox33/rick/http_framework/framework/middleware"
	"time"
)

func registerRouter(core *framework.Core) {
	core.Use(middleware.RecoveryMiddleware())
	core.Get("/foo", FooControllerHandler)
	core.Get("/user/login", UserLoginHandler)
	subjectApi := core.Group("/subject")
	{
		subjectApi.Delete("/:id", middleware.TimeoutMiddleware(time.Second), SubjectDelHandler)
		subjectApi.Put("/:id", SubjectUpdateHandler)
		subjectApi.Get("/:id", SubjectGetHandler)
		subjectApi.Get("/list/all", SubjectListHandler)
	}
}
