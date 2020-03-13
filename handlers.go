package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const StatusOnline = "online"

func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, StatusOnline)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		handleJSONDecodeError(&w, RegisterName, err)
		return
	}

	log.Infof("%s: %+v\n", RegisterName, request)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		handleJSONDecodeError(&w, LoginName, err)
		return
	}

	log.Infof("%s: %+v\n", LoginName, request)
}