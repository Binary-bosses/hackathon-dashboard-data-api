package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Binary-bosses/hackathon-dashboard-data-api/util"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// Server -
type Server interface {
	Start(string) error
	Stop()
}

// StartServer creates and starts a new server on the specified path
func StartServer(path string) error {
	srv, err := newServer()
	if err != nil {
		return err
	}
	return fasthttp.ListenAndServe(path, srv.router.Handler)
}

// NewServer creates a new server, but doesn't start it
func NewServer() (Server, error) {
	return newServer()
}

func (s *server) Start(path string) error {
	return s.srv.ListenAndServe(path)
}

func (s *server) Stop() {
	s.srv.Shutdown()
}

type server struct {
	srv           *fasthttp.Server
	router        *fasthttprouter.Router
	defaultClient *http.Client
}

func newServer() (*server, error) {

	client := &http.Client{Timeout: time.Second * 30}

	srv := &server{
		router:        fasthttprouter.New(),
		defaultClient: client,
	}
	srv.setupRoutes()
	srv.srv = &fasthttp.Server{
		Handler:        srv.router.Handler,
		ReadBufferSize: 8192,
	}

	return srv, nil
}

// APIResponse is wrapped around responses from the API
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

// BasicResponse is a simple wrapper for an exception response
func BasicResponse(status int, msg string, ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode(status)
	util.SetJSONBody(ctx, APIResponse{Status: status, Message: msg})
}

func (s *server) echo() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		resp := ctx.UserValue("echo")
		if resp == "" {
			resp = "hello"
		}
		ctx.SetStatusCode(200)
		fmt.Fprintf(ctx, "%s", resp)
	}
}
