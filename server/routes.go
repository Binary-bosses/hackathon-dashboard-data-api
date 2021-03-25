package server

import (
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Binary-bosses/hackathon-dashboard-data-api/util"
	"github.com/valyala/fasthttp"
)

var funcCleaner = regexp.MustCompile(`^(.*?)\.func\d+$`)

type route struct {
	httpMethod     string
	path           string
	requestHandler fasthttp.RequestHandler
	hideLog        bool
}

func (s *server) routes() []route {
	return []route{{"GET", "/api/v1/echo/:echo", s.echo(), false}}
}

// setupRoutes configures all the route info, automatically adding logging
func (s *server) setupRoutes() {

	s.router.NotFound = func(ctx *fasthttp.RequestCtx) {
		apiResp := APIResponse{
			Status: 404,
			Data:   "Path does not exist",
		}
		ctx.Response.SetStatusCode(404)
		util.SetJSONBody(ctx, apiResp)
	}

	for _, r := range s.routes() {
		currFunc := r.requestHandler

		if !r.hideLog {
			currFunc = s.logWrapper(currFunc)
		}

		s.router.Handle(r.httpMethod, r.path, currFunc)
		log.Printf("Mapped route %s %s to %s\n", r.httpMethod, r.path, getFunctionName(r.requestHandler))
	}
}

// Wraps a logger around the response that gives the name of the responding function,
// http method, and path. It also tells how long it took to complete, and its status
func (s *server) logWrapper(f fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Printf("Received request ID %d: %s %s", ctx.ConnID(), ctx.Method(), ctx.RequestURI())

		start := time.Now()
		f(ctx)
		end := time.Since(start)

		log.Printf("Finished request ID %d: %s %s - HTTP %d %s (%.3fs)", ctx.ConnID(), ctx.Method(), ctx.RequestURI(), ctx.Response.StatusCode(), http.StatusText(ctx.Response.StatusCode()), end.Seconds())
	}
}

func replaceSubstring(s, substr, repl string) string {
	if i := strings.Index(s, substr); i > -1 {
		s = s[:i] + repl + s[i+len(substr):]
	}
	return s
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

// Returns a human friendly name of the specified function
func getFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	name = funcCleaner.ReplaceAllString(name, "$1")
	name = replaceSubstring(name, ".(*server).", ".")
	return name
}
