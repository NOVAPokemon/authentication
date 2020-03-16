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

	log.Infof("Starting AUTHENTICATION server in port %d...\n", Port)
	log.Fatal(http.ListenAndServe(addr, router))
}
