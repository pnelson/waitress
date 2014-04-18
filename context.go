package waitress

import (
	"encoding/json"
	"net/http"

	"github.com/pnelson/waitress/router"
)

type Context struct {
	Response *ResponseWriter
	Request  *http.Request
	adapter  *router.Adapter
}

func NewContext(w *ResponseWriter, r *http.Request, adapter *router.Adapter) *Context {
	return &Context{
		Response: w,
		Request:  r,
		adapter:  adapter,
	}
}

func (ctx *Context) Abort(code int) http.Handler {
	switch code {
	case 400:
		return BadRequest()
	case 401:
		return Unauthorized()
	case 403:
		return Forbidden()
	case 404:
		return NotFound() // TODO: use router NotFound
	case 406:
		return NotAcceptable()
	case 408:
		return RequestTimeout()
	case 409:
		return Conflict()
	case 410:
		return Gone()
	case 411:
		return LengthRequired()
	case 412:
		return PreconditionFailed()
	case 413:
		return RequestEntityTooLarge()
	case 414:
		return RequestURITooLong()
	case 415:
		return UnsupportedMediaType()
	case 416:
		return RequestedRangeNotSatisfiable()
	case 417:
		return ExpectationFailed()
	case 422:
		return UnprocessableEntity()
	case 429:
		return TooManyRequests()
	case 501:
		return NotImplemented()
	case 502:
		return BadGateway()
	case 503:
		return ServiceUnavailable()
	case 504:
		return GatewayTimeout()
	}
	return InternalServerError() // TODO: use router InternalServerError
}

func (ctx *Context) Build(method, name string) *router.Builder {
	return ctx.adapter.Build(method, name)
}

func (ctx *Context) ByteHandler(b []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})
}

func (ctx *Context) DecodeJSON(i interface{}) error {
	decoder := json.NewDecoder(ctx.Request.Body)
	return decoder.Decode(i)
}

func (ctx *Context) EncodeJSON(i interface{}) ([]byte, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (ctx *Context) WriteJSON(b []byte) http.Handler {
	ctx.Header("Content-Type", "application/json")
	return ctx.ByteHandler(b)
}

func (ctx *Context) Header(key, value string) {
	ctx.Response.Header().Set(key, value)
}

func (ctx *Context) Redirect(builder *router.Builder) http.Handler {
	return ctx.RedirectWithCode(builder, 303)
}

func (ctx *Context) RedirectWithCode(builder *router.Builder, code int) http.Handler {
	path, ok := builder.Full()
	if !ok {
		return ctx.Abort(500)
	}

	return ctx.RedirectToWithCode(path, code)
}

func (ctx *Context) RedirectTo(path string) http.Handler {
	return ctx.RedirectToWithCode(path, 303)
}

func (ctx *Context) RedirectToWithCode(path string, code int) http.Handler {
	return RedirectToWithCode(path, code)
}

func (ctx *Context) Status(code int) {
	ctx.Response.status = code
}
