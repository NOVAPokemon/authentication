package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/cookies"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
	userdb "github.com/NOVAPokemon/utils/database/user"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Constants
const StatusOnline = "online"

// Indicates if the server is online.
func Status(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, StatusOnline)
}

// Registers a user. Expects a JSON with username and password in the body.
func Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.HandleJSONDecodeError(&w, RegisterName, err)
		return
	}

	trainerToAdd := utils.Trainer{
		Username: request.Username,
		Items:    map[string]utils.Item{},
		Pokemons: map[string]utils.Pokemon{},
		Stats: utils.TrainerStats{
			Level: 0,
			Coins: 0,
		},
	}

	userToAdd := utils.User{
		Username:     request.Username,
		PasswordHash: hashPassword([]byte(request.Password)),
	}

	err, id := userdb.AddUser(&userToAdd)

	if err != nil {
		return
	}

	trainerId, err := trainerdb.AddTrainer(trainerToAdd)

	if err != nil {
		return
	}

	log.Infof("%s: %s %s %s\n", RegisterName, request.Username, id, trainerId)
}

// Logs in a user. Expects a JSON with username and password in the body.
func Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.HandleJSONDecodeError(&w, LoginName, err)
		return
	}

	err, user := userdb.GetUserByUsername(request.Username)

	if err != nil {
		return
	}

	if !verifyPassword([]byte(request.Password), user.PasswordHash) {
		utils.HandleWrongPasswordError(&w, LoginName, request.Username, request.Password)
		return
	}

	err = cookies.SetAuthToken(request.Username, LoginName, &w)

	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("%s: %s\n", LoginName, request.Username)
}

// Endpoint to refresh the token. Expects the user to already have a token.
func Refresh(w http.ResponseWriter, r *http.Request) {
	claims, err := cookies.ExtractAndVerifyAuthToken(&w, r, RefreshName)

	if err != nil {
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 2*time.Minute {
		utils.HandleToSoonToRefreshError(&w, RefreshName)
		return
	}

	err = cookies.SetAuthToken(claims.Username, RefreshName, &w)

	if err != nil {
		return
	}

	log.Infof("%s: %s\n", RefreshName, claims.Username)
}

func hashPassword(password []byte) (passwordHash []byte) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Error(err)
	}

	return hash
}

func verifyPassword(password, expectedHash []byte) bool {
	err := bcrypt.CompareHashAndPassword(expectedHash, password)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
