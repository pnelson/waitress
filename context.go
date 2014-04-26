package waitress

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pnelson/waitress/router"
)

// Context exists as a collection of endpoint helper methods. You have full
// access to the http.ResponseWriter and http.Request as you would with the
// standard library.
type Context struct {
	Response *ResponseWriter
	Request  *http.Request
	adapter  *router.Adapter
}

// NewContext returns a NewContext bound to the provided parameters.
func NewContext(w *ResponseWriter, r *http.Request, adapter *router.Adapter) *Context {
	return &Context{
		Response: w,
		Request:  r,
		adapter:  adapter,
	}
}

// Abort returns an http.Handler for HTTP status codes greater than 400. Some
// error codes have not been implemented.
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

// Build returns a Builder from the adapter preconfigured for the request.
func (ctx *Context) Build(method, name string) *router.Builder {
	return ctx.adapter.Build(method, name)
}

// ByteHandler returns an http.Handler that writes the provided bytes.
func (ctx *Context) ByteHandler(b []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})
}

// DecodeJSON decodes the request body into a struct. Provide a pointer to a
// struct as you would with the encoding/json package.
func (ctx *Context) DecodeJSON(i interface{}) error {
	decoder := json.NewDecoder(ctx.Request.Body)
	return decoder.Decode(i)
}

// EncodeJSON marshals the provided structure into a byte slice.
func (ctx *Context) EncodeJSON(i interface{}) ([]byte, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// QueryInt converts a request query parameter to an integer. This is a
// convenience method for the common scenario of accepting query parameters to
// modify result sets, like pagination.
func (ctx *Context) QueryInt(key string, def int64, base, size int) (int64, error) {
	query := ctx.Request.URL.Query()
	value := query.Get(key)
	if value == "" {
		return def, nil
	}

	rv, err := strconv.ParseInt(value, base, size)
	if err != nil {
		return def, err
	}

	return rv, nil
}

// QueryInt64 is like QueryInt but specifies base 10 and 64 bit result.
func (ctx *Context) QueryInt64(key string, def int64) (int64, error) {
	return ctx.QueryInt(key, def, 10, 64)
}

// WriteJSON sets the response content type to application/json and returns an
// http.Handler that will write the bytes to the response.
func (ctx *Context) WriteJSON(b []byte) http.Handler {
	ctx.Header("Content-Type", "application/json")
	return ctx.ByteHandler(b)
}

// Header is a shortcut to setting a response header.
func (ctx *Context) Header(key, value string) {
	ctx.Response.Header().Set(key, value)
}

// Redirect will perform an HTTP 303 redirect to the endpoint built by the
// provided Builder.
func (ctx *Context) Redirect(builder *router.Builder) http.Handler {
	return ctx.RedirectWithCode(builder, 303)
}

// RedirectWithCode is like Redirect but you provide the HTTP status code to be
// used for redirection.
func (ctx *Context) RedirectWithCode(builder *router.Builder, code int) http.Handler {
	url, ok := builder.Build()
	if !ok {
		return ctx.Abort(500)
	}

	return ctx.RedirectToWithCode(url.String(), code)
}

// RedirectTo will perform an HTTP 303 redirect to the provided endpoint. You
// should provide an absolute URI here.
func (ctx *Context) RedirectTo(path string) http.Handler {
	return ctx.RedirectToWithCode(path, 303)
}

// RedirectToWithCode is like RedirectTo but you provide the HTTP status code
// to be used for redirection.
func (ctx *Context) RedirectToWithCode(path string, code int) http.Handler {
	return RedirectToWithCode(path, code)
}

// Status configures the response status code to write to the header. The
// response header will not be written by calling this method.
func (ctx *Context) Status(code int) {
	ctx.Response.status = code
}
