package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/clients"
	userdb "github.com/NOVAPokemon/utils/database/user"
	"github.com/NOVAPokemon/utils/pokemons"
	"github.com/NOVAPokemon/utils/tokens"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Constants
const StatusOnline = "online"

var httpClient = &http.Client{}

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

	log.Infof("Register request for: %s", request.Username)

	exists, err := userdb.CheckIfUserExists(request.Username)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred registering user", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "User Already exists", http.StatusConflict)
		return
	}

	userToAdd := utils.User{
		Username:     request.Username,
		PasswordHash: hashPassword([]byte(request.Password)),
	}

	err, id := userdb.AddUser(&userToAdd)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred registering user", http.StatusInternalServerError)
		return
	}

	trainerToAdd := utils.Trainer{
		Username: request.Username,
		Pokemons: generateStarterPokemons(6),
		Stats: utils.TrainerStats{
			Level: 0,
			Coins: 200,
		},
	}

	trainersClient := clients.NewTrainersClient(fmt.Sprintf("%s:%d", host, utils.TrainersPort), httpClient)
	_, err = trainersClient.AddTrainer(trainerToAdd)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred registering trainer", http.StatusInternalServerError)
		return
	}

	log.Infof("%s: %s %s %s\n", RegisterName, request.Username, id, id)
}

// Logs in a user. Expects a JSON with username and password in the body.
func Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.HandleJSONDecodeError(&w, LoginName, err)
		return
	}

	log.Infof("Login request for: %s", request.Username)

	err, user := userdb.GetUserByUsername(request.Username)

	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	if !verifyPassword([]byte(request.Password), user.PasswordHash) {
		utils.HandleWrongPasswordError(&w, LoginName, request.Username, request.Password)
		return
	}

	tokens.AddAuthToken(request.Username, w.Header())
	log.Infof("%s: %s\n", LoginName, request.Username)
}

// Endpoint to refresh the token. Expects the user to already have a token.
func Refresh(w http.ResponseWriter, r *http.Request) {
	claims, err := tokens.ExtractAndVerifyAuthToken(r.Header)

	if err != nil {
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 2*time.Minute {
		utils.HandleToSoonToRefreshError(&w, RefreshName)
		return
	}

	tokens.AddAuthToken(claims.Username, w.Header())
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

/*
  "max_level": 100.0,
  "max_hp": 500.0,
  "max_damage": 500.0,
  "StdHpDeviation": 25.0,
  "StdDamageDeviation": 25.0
*/

func generateStarterPokemons(pokemonNr int) map[string]pokemons.Pokemon { //TODO only for testing

	toReturn := make(map[string]pokemons.Pokemon, pokemonNr)

	for i := 0; i < pokemonNr; i++ {
		newPokemon := pokemons.GetOneWildPokemon(
			5,
			500,
			500,
			25,
			300,
			fmt.Sprintf("starter-%d", i))
		toReturn[newPokemon.Id.Hex()] = *newPokemon
	}
	return toReturn

}
