package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const host = "localhost"
const port = 8080

func main() {
	router := NewRouter()
	addr := fmt.Sprintf("%s:%d", host, port)

	log.Info("Starting AUTHENTICATION server...")
	log.Fatal(http.ListenAndServe(addr, router))
}
