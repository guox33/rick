package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func NewContext(request *http.Request, response http.ResponseWriter) *Context {
	return &Context{
		request:    request,
		response:   response,
		writeMutex: &sync.Mutex{},
		hasTimeout: false,
		index:      -1,
	}
}

type Context struct {
	request  *http.Request
	response http.ResponseWriter

	// Write lock
	writeMutex *sync.Mutex
	// Flag which represent timeout or not
	hasTimeout     bool
	requestTimeout *time.Duration

	handlerChain ControlHandlerChain
	index        int8
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlerChain)) {
		c.handlerChain[c.index](c)
		c.index++
	}
}

func (c *Context) WriteMux() *sync.Mutex {
	return c.writeMutex
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Response() http.ResponseWriter {
	return c.response
}

func (c *Context) SetHandler() {

}

func (c *Context) HasTimeout() bool {
	return c.hasTimeout
}

func (c *Context) SetHasTimeout() {
	c.hasTimeout = true
	return
}

func (c *Context) RequestTimeout() *time.Duration {
	return c.requestTimeout
}

func (c *Context) SetRequestTimeout(t time.Duration) {
	c.requestTimeout = &t
	return
}

func (c *Context) BaseContext() context.Context {
	return c.request.Context()
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.request.Context().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.request.Context().Done()
}

func (c *Context) Err() error {
	return c.request.Context().Err()
}

func (c *Context) Value(key any) any {
	return c.request.Context().Value(key)
}

func (c *Context) QueryInt(key string, def int) int {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			val := params[len(params)-1]
			if v, err := strconv.Atoi(val); err == nil {
				return v
			} else {
				return def
			}
		}
	}
	return 0
}

func (c *Context) QueryString(key string) string {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			return params[len(params)-1]
		}
	}
	return ""
}

func (c *Context) QueryArray(key string) []string {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		return params
	}
	return []string{}
}

func (c *Context) QueryAll() map[string][]string {
	if c.request != nil {
		return c.request.URL.Query()
	}
	return map[string][]string{}
}

func (c *Context) FormInt(key string, def int) int {
	args := c.FormAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			val := params[len(params)-1]
			if v, err := strconv.Atoi(val); err == nil {
				return v
			} else {
				return def
			}
		}
	}
	return 0
}

func (c *Context) FormString(key string) string {
	args := c.FormAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			return params[len(params)-1]
		}
	}
	return ""
}

func (c *Context) FormArray(key string) []string {
	args := c.FormAll()
	if params, ok := args[key]; ok {
		return params
	}
	return []string{}
}

func (c *Context) FormAll() map[string][]string {
	if c.request != nil {
		return c.request.PostForm
	}
	return map[string][]string{}
}

func (c *Context) BindJson(data interface{}) error {
	if c.request == nil {
		return errors.New("request is nil")
	}

	buf, err := io.ReadAll(c.request.Body)
	if err != nil {
		return err
	}
	c.request.Body = io.NopCloser(bytes.NewBuffer(buf))

	err = json.Unmarshal(buf, &data)
	return err
}

func (c *Context) Json(status int, data interface{}) error {
	if c.hasTimeout {
		return nil
	}

	log.Default().Println("text: ", data)

	c.response.WriteHeader(status)
	c.response.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(data)
	if err != nil {
		c.response.WriteHeader(500)
		return err
	}
	_, err = c.response.Write(body)
	if err != nil {
		c.response.WriteHeader(500)
		return err
	}
	return nil
}

func (c *Context) Html(status int, data interface{}, template string) error {
	if c.hasTimeout {
		return nil
	}

	return nil
}

func (c *Context) Text(status int, data string) error {
	if c.hasTimeout {
		return nil
	}

	if c.request == nil {
		return errors.New("request is nil")
	}

	log.Default().Println("text: " + data)

	c.response.WriteHeader(status)
	c.response.Header().Set("Content-Type", "text/plain")
	_, err := c.response.Write([]byte(data))
	return err
}

func (c *Context) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(404)
	_, _ = response.Write([]byte("not found not found"))
	context.Background()
}
