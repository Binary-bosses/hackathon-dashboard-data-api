package util

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// SetJSONBody sets the response body as json with the given payload
func SetJSONBody(ctx *fasthttp.RequestCtx, i interface{}) {
	res, err := json.Marshal(i)
	ctx.Response.Header.Set("Content-type", "application/json")
	if err == nil {
		ctx.Response.SetBody(res)
	}
}
