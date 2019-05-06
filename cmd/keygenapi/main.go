package main

import (
	"fmt"
	"flag"
	"keygen/keygenapi"
	"keygen"
	"log"
	"net/http"
)

var (
	port = flag.Int("p", 8080, "port to listen http requests")
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	flag.Parse()

	s := keygen.NewStore()
	kg := keygen.New(s)
	a := keygenapi.New(kg)

	http.HandleFunc("/healthcheck", healthCheck)
	http.Handle("/", a)

	log.Printf("Server started at port %d", *port)

	po := fmt.Sprintf(":%d", *port)
	log.Fatal(
		http.ListenAndServe(po, nil),
	)
}
