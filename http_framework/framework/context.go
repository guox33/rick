package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
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

	fullPath string
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

func (c *Context) QueryInt(key string, def int) (int, bool) {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			if v, err := cast.ToIntE(params[0]); err == nil {
				return v, true
			}
		}
	}
	return def, false
}

func (c *Context) QueryInt64(key string, def int64) (int64, bool) {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			if v, err := cast.ToInt64E(params[0]); err == nil {
				return v, true
			}
		}
	}
	return def, false
}

func (c *Context) QueryFloat64(key string, def float64) (float64, bool) {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			if v, err := cast.ToFloat64E(params[0]); err == nil {
				return v, true
			}
		}
	}
	return def, false
}

func (c *Context) QueryFloat32(key string, def float32) (float32, bool) {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			if v, err := cast.ToFloat32E(params[0]); err == nil {
				return v, true
			}
		}
	}
	return def, false
}

func (c *Context) QueryBool(key string, def bool) (bool, bool) {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			if v, err := cast.ToBoolE(params[0]); err == nil {
				return v, true
			}
		}
	}
	return def, false
}

func (c *Context) QueryString(key, def string) (string, bool) {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			return params[0], true
		}
	}
	return def, false
}

func (c *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		return params, true
	}
	return def, false
}

func (c *Context) Query(key string) interface{} {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		return params
	}
	return nil
}

func (c *Context) QueryAll() map[string][]string {
	if c.request != nil {
		return c.request.URL.Query()
	}
	return map[string][]string{}
}

func (c *Context) ParamInt(key string, def int) (int, bool) {
	args := c.ParamAll()
	params, ok := args[key]
	if !ok {
		return def, false
	}
	if v, err := cast.ToIntE(params); err == nil {
		return v, true
	}
	return def, false
}

func (c *Context) ParamInt64(key string, def int64) (int64, bool) {
	args := c.ParamAll()
	params, ok := args[key]
	if !ok {
		return def, false
	}
	if v, err := cast.ToInt64E(params); err == nil {
		return v, true
	}
	return def, false
}

func (c *Context) ParamFloat64(key string, def float64) (float64, bool) {
	args := c.ParamAll()
	params, ok := args[key]
	if !ok {
		return def, false
	}
	if v, err := cast.ToFloat64E(params); err == nil {
		return v, true
	}
	return def, false
}

func (c *Context) ParamFloat32(key string, def float32) (float32, bool) {
	args := c.ParamAll()
	params, ok := args[key]
	if !ok {
		return def, false
	}
	if v, err := cast.ToFloat32E(params); err == nil {
		return v, true
	}
	return def, false
}

func (c *Context) ParamBool(key string, def bool) (bool, bool) {
	args := c.ParamAll()
	params, ok := args[key]
	if !ok {
		return def, false
	}
	if v, err := cast.ToBoolE(params); err == nil {
		return v, true
	}
	return def, false
}

func (c *Context) ParamString(key, def string) (string, bool) {
	args := c.ParamAll()
	params, ok := args[key]
	if !ok {
		return def, false
	}
	return params, true
}

func (c *Context) ParamStringSlice(key string, def []string) ([]string, bool) {
	args := c.QueryAll()
	if params, ok := args[key]; ok {
		return params, true
	}
	return def, false
}

func (c *Context) Param(key string) interface{} {
	args := c.ParamAll()
	params, _ := args[key]
	return params
}

func (c *Context) ParamAll() map[string]string {
	if c.request == nil {
		return map[string]string{}
	}

	params := make(map[string]string)

	segments := strings.Split(c.fullPath, "/")
	uriSegments := strings.Split(c.request.RequestURI, "/")
	for i := range segments {
		if !strings.HasPrefix(segments[i], ":") {
			continue
		}
		params[strings.TrimPrefix(segments[i], ":")] = uriSegments[i]
	}

	return params
}

func (c *Context) FormInt64(key string, def int64) (int64, bool) {
	args := c.FormAll()
	params, ok := args[key]
	if !ok || len(params) == 0 {
		return def, false
	}
	if val, err := cast.ToInt64E(params[0]); err == nil {
		return val, true
	}
	return def, false
}

func (c *Context) FormInt(key string, def int) (int, bool) {
	args := c.FormAll()
	params, ok := args[key]
	if !ok || len(params) == 0 {
		return def, false
	}
	if val, err := cast.ToIntE(params[0]); err == nil {
		return val, true
	}
	return def, false
}

func (c *Context) FormFloat64(key string, def float64) (float64, bool) {
	args := c.FormAll()
	params, ok := args[key]
	if !ok || len(params) == 0 {
		return def, false
	}
	if val, err := cast.ToFloat64E(params[0]); err == nil {
		return val, true
	}
	return def, false
}

func (c *Context) FormFloat32(key string, def float32) (float32, bool) {
	args := c.FormAll()
	params, ok := args[key]
	if !ok || len(params) == 0 {
		return def, false
	}
	if val, err := cast.ToFloat32E(params[0]); err == nil {
		return val, true
	}
	return def, false
}

func (c *Context) FormBool(key string, def bool) (bool, bool) {
	args := c.FormAll()
	params, ok := args[key]
	if !ok || len(params) == 0 {
		return def, false
	}
	if val, err := cast.ToBoolE(params[0]); err == nil {
		return val, true
	}
	return def, false
}

func (c *Context) FormString(key string, def string) (string, bool) {
	args := c.FormAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			return params[0], true
		}
	}
	return def, false
}

func (c *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	args := c.FormAll()
	if params, ok := args[key]; ok {
		if len(params) > 0 {
			return params, true
		}
	}
	return def, false
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if err := c.request.ParseMultipartForm(1024000); err != nil {
		return nil, err
	}
	f, fh, err := c.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	_ = f.Close()
	return fh, nil
}

func (c *Context) FormAll() map[string][]string {
	if c.request != nil {
		return map[string][]string{}
	}
	if err := c.request.ParseMultipartForm(1024000); err != nil {
		return map[string][]string{}
	}
	return c.request.PostForm
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

func (c *Context) BindXml(data interface{}) error {
	if c.request == nil {
		return errors.New("request is nil")
	}

	buf, err := io.ReadAll(c.request.Body)
	if err != nil {
		return err
	}
	c.request.Body = io.NopCloser(bytes.NewBuffer(buf))

	err = xml.Unmarshal(buf, &data)
	return err
}

func (c *Context) GetRawData() ([]byte, error) {
	if c.request == nil {
		return nil, errors.New("request is nil")
	}

	buf, err := io.ReadAll(c.request.Body)
	if err != nil {
		return nil, err
	}
	c.request.Body = io.NopCloser(bytes.NewBuffer(buf))

	return buf, nil
}

func (c *Context) Uri() string {
	if c.request == nil {
		return ""
	}
	return c.request.RequestURI
}

func (c *Context) Method() string {
	if c.request == nil {
		return ""
	}
	return c.request.Method
}

func (c *Context) Host() string {
	if c.request == nil {
		return ""
	}
	return c.request.Host
}

func (c *Context) ClientIp() string {
	if c.request == nil {
		return ""
	}
	return ""
}

func (c *Context) Json(data interface{}) IResponse {
	if c.hasTimeout {
		return nil
	}

	log.Default().Println("text: ", data)

	c.response.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(data)
	if err != nil {
		return c
	}
	_, err = c.response.Write(body)
	if err != nil {
		return c
	}
	return nil
}

func (c *Context) Jsonp(obj interface{}) IResponse {
	callbackFunc, _ := c.QueryString("callback", "callback_function")
	c.SetHeader("Content-Type", "application/javascript")
	callback := template.JSEscapeString(callbackFunc)

	_, err := c.response.Write([]byte(callback))
	if err != nil {
		return c
	}
	_, err = c.response.Write([]byte("("))
	if err != nil {
		return c
	}
	ret, err := json.Marshal(obj)
	if err != nil {
		return c
	}
	_, err = c.response.Write(ret)
	if err != nil {
		return c
	}
	_, err = c.response.Write([]byte(")"))
	return c
}

func (c *Context) Xml(data interface{}) IResponse {
	return c
}

func (c *Context) Html(file string, data interface{}) IResponse {
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return c
	}

	if err := t.Execute(c.response, data); err != nil {
		return c
	}

	return c
}

func (c *Context) Text(data string) IResponse {
	if c.hasTimeout {
		return nil
	}

	if c.response == nil {
		return c
	}

	log.Default().Println("text: " + data)

	c.response.Header().Set("Content-Type", "text/plain")
	_, _ = c.response.Write([]byte(data))
	return c
}

func (c *Context) Redirect(path string) IResponse {
	return c
}

func (c *Context) SetHeader(key string, val string) IResponse {
	c.response.Header().Set(key, val)
	return c
}

func (c *Context) SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	http.SetCookie(c.response, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: 0,
	})
	return c
}

func (c *Context) SetStatus(code int) IResponse {
	if c.response == nil {
		return c
	}
	c.response.WriteHeader(code)
	return c
}

func (c *Context) SetOKStatus() IResponse {
	if c.response == nil {
		return c
	}
	c.response.WriteHeader(http.StatusOK)
	return c
}

func (c *Context) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(404)
	_, _ = response.Write([]byte("not found not found"))
	context.Background()
}
