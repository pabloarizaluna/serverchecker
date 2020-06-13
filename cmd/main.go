package main

import (
	"log"

	"github.com/pabloarizaluna/serverchecker/cockroach"
	"github.com/pabloarizaluna/serverchecker/web"
	"github.com/valyala/fasthttp"
)

func main() {
	store, err := cockroach.NewStore(
		"postgresql://craig@localhost:26257/checker?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.craig.key&sslcert=certs/client.craig.crt")
	if err != nil {
		log.Fatal(err)
	}
	h := web.NewHandler(store)
	log.Fatal(fasthttp.ListenAndServe(":3000", h.Handler))
}
