package context

import (
	"net/http"

	"github.com/shishir1290/goduck/runtime/binder"
	"github.com/shishir1290/goduck/runtime/request"
	"github.com/shishir1290/goduck/runtime/response"
)

type Context struct {
	Request  *request.Request
	Response *response.Response

	params map[string]string
}

func New(
	w http.ResponseWriter,
	r *http.Request,
) *Context {

	recorder := response.NewRecorder(w)

	return &Context{

		Request: request.New(r),

		Response: response.New(recorder),

		params: make(map[string]string),
	}
}

func (ctx *Context) String(text string) {
	ctx.Response.String(text)
}

func (ctx *Context) JSON(v any) error {
	return ctx.Response.JSON(v)
}

func (ctx *Context) Status(code int) {
	ctx.Response.Status(code)
}

func (ctx *Context) Param(name string) string {
	return ctx.params[name]
}

func (ctx *Context) SetParam(name, value string) {
	ctx.params[name] = value
}

func (ctx *Context) Bind(v any) error {

	return binder.Bind(
		ctx.Request.Raw(),
		v,
	)

}

func (ctx *Context) StatusCode() int {

	return ctx.Response.StatusCode()
}

func (ctx *Context) ResponseSize() int {

	return ctx.Response.Size()
}