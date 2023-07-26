package framework

type IResponse interface {
	Json(obj interface{}) IResponse
	Jsonp(obj interface{}) IResponse
	Xml(obj interface{}) IResponse
	Html(file string, obj interface{}) IResponse
	Text(obj string) IResponse
	Redirect(path string) IResponse

	SetHeader(key string, val string) IResponse
	SetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse
	SetStatus(code int) IResponse
	SetOKStatus() IResponse
}
