package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handleJSONDecodeError(w *http.ResponseWriter, caller string, err error) {
	log.Warnf("Error decoding body from %s request:\n", caller)
	log.Warn(err)
	(*w).WriteHeader(400)
}
