package web

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pabloarizaluna/serverchecker"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	*fasthttprouter.Router

	store serverchecker.Store
}

func NewHandler() *Handler {
	h := &Handler{Router: fasthttprouter.New()}
	h.GET("/check/:domain", h.checkDomain)
	h.GET("/domain", h.domains)

	return h
}

func (s *Handler) checkDomain(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Checking...")
	request := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", ctx.UserValue("domain"))
	resp, err := http.Get(request)
	if err != nil {
		fmt.Fprint(ctx, err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(ctx, err.Error())
	}

	fmt.Fprint(ctx, string(body))
}

func (s *Handler) domains(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Loading items...")
}
