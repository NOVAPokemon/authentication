package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handleJSONDecodeError(w *http.ResponseWriter, caller string, err error) {
	log.Warnf("Error decoding body from %s request:\n", caller)
	log.Warn(err)
	(*w).WriteHeader(http.StatusBadRequest)
}

func handleWrongPasswordError(w *http.ResponseWriter, caller, username, password string) {
	log.Warnf("Error wrong password in %s request:\n", caller)
	log.Warn(fmt.Sprintf("%s : %s", username, password))
	(*w).WriteHeader(http.StatusUnauthorized)
}

func handleJWTSigningError(w *http.ResponseWriter, caller string, err error) {
	log.Warnf("Error signing jwt in %s request:\n", caller)
	log.Warn(err)
	(*w).WriteHeader(http.StatusInternalServerError)
}

func handleCookieError(w *http.ResponseWriter, caller string, err error) {
	if err == http.ErrNoCookie {
		log.Warnf("Error no cookie in %s request:\n", caller)
		log.Warn(err)
		(*w).WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Warnf("Error other in %s request:\n", caller)
	log.Warn(err)
	(*w).WriteHeader(http.StatusBadRequest)
	return
}

func handleJWTVerifyingError(w *http.ResponseWriter, caller string, err error) {
	if err == jwt.ErrSignatureInvalid {
		log.Warnf("Error invalid signature in %s request:\n", caller)
		log.Warn(err)
		(*w).WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Warnf("Error other in %s request:\n", caller)
	log.Warn(err)
	(*w).WriteHeader(http.StatusBadRequest)
	return
}
