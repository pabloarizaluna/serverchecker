package web

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// Server represents the services
type Server struct {
	R *fasthttprouter.Router
}

// NewServer return a server configured
func NewServer() *Server {
	s := &Server{fasthttprouter.New()}
	s.R.GET("/check/:domain", s.checkDomain)
	s.R.GET("/domain", s.domains)

	return s
}

func (s *Server) checkDomain(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Checking...")
	request := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", ctx.UserValue("domain"))
	resp, err := http.Get(request)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Fprint(ctx, string(body))
}

func (s *Server) domains(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Loading items...")
}
