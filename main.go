package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const host = "localhost"
const Port = 8001

func main() {
	router := utils.NewRouter(routes)
	addr := fmt.Sprintf("%s:%d", host, Port)

	log.Info("Starting AUTHENTICATION server...")
	log.Fatal(http.ListenAndServe(addr, router))
}
