package main

import (
	"log"

	"github.com/pabloarizaluna/server-checker/web"
	"github.com/valyala/fasthttp"
)

func main() {
	s := web.NewServer()
	log.Fatal(fasthttp.ListenAndServe(":8080", s.R.Handler))
}
