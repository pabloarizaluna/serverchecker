package main

import (
	"log"

	"github.com/pabloarizaluna/serverchecker/cockroach"
	"github.com/pabloarizaluna/serverchecker/web"
	"github.com/valyala/fasthttp"
)

const (
	dataSource           = "postgresql://craig@localhost:26257/checker?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.craig.key&sslcert=certs/client.craig.crt"
	corsAllowMethods     = "GET"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "false"
)

func cors(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.SetContentType("application/json")

		next(ctx)
	}
}

func main() {
	store, err := cockroach.NewStore(dataSource)
	if err != nil {
		log.Fatal(err)
	}
	h := web.NewHandler(store)
	if err := fasthttp.ListenAndServe(":3000", cors(h.Handler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
