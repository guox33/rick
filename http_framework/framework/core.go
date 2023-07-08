package framework

import "net/http"

type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(404)
	_, _ = response.Write([]byte("not found not found"))
}
