package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const host = utils.ServeHost
const port = utils.AuthenticationPort

func main() {
	router := utils.NewRouter(routes)
	addr := fmt.Sprintf("%s:%d", host, port)

	log.Infof("Starting AUTHENTICATION server in port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, router))
}
