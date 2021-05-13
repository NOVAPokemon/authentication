package main

import (
	"encoding/json"
	"fmt"
	"os"

	http "github.com/bruno-anjos/archimedesHTTPClient"

	originalHTTP "net/http"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/clients"
	userdb "github.com/NOVAPokemon/utils/database/user"
	"github.com/NOVAPokemon/utils/pokemons"
	"github.com/NOVAPokemon/utils/tokens"
	"github.com/NOVAPokemon/utils/websockets"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	httpClient = &http.Client{
		Client: originalHTTP.Client{
			Timeout:   websockets.Timeout,
			Transport: clients.NewTransport(),
		},
	}
	serverName   string
	commsManager websockets.CommunicationManager
)

func init() {
	if aux, exists := os.LookupEnv(utils.HostnameEnvVar); exists {
		serverName = aux
	} else {
		log.Fatal("Could not load server name")
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	var request registerRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRegisterHandlerError(err), http.StatusBadRequest)
		return
	}

	log.Infof("Register request for: %s", request.Username)

	exists, err := userdb.CheckIfUserExists(request.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRegisterHandlerError(err), http.StatusInternalServerError)
		return
	}

	if exists {
		err = wrapRegisterHandlerError(newRegisterConflictError(request.Username))
		utils.LogAndSendHTTPError(&w, err, http.StatusConflict)
		return
	}

	hash, err := hashPassword([]byte(request.Password))
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRegisterHandlerError(err), http.StatusInternalServerError)
		return
	}

	userToAdd := utils.User{
		Username:     request.Username,
		PasswordHash: hash,
	}

	id, err := userdb.AddUser(&userToAdd)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRegisterHandlerError(err), http.StatusInternalServerError)
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

	trainersClient := clients.NewTrainersClient(httpClient, commsManager)
	_, err = trainersClient.AddTrainer(trainerToAdd)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRegisterHandlerError(err), http.StatusInternalServerError)
		return
	}

	log.Infof("%s: %s %s %s\n", registerName, request.Username, id, id)
}

func login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapLoginHandlerError(err), http.StatusBadRequest)
		return
	}

	log.Infof("Login request for: %s", request.Username)

	user, err := userdb.GetUserByUsername(request.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapLoginHandlerError(err), http.StatusNotFound)
		return
	}

	if !verifyPassword([]byte(request.Password), user.PasswordHash) {
		err = wrapLoginHandlerError(newWrongPasswordError(request.Username))
		utils.LogAndSendHTTPError(&w, err, http.StatusBadRequest)
		return
	}

	tokens.AddAuthToken(request.Username, w.Header())
	log.Infof("%s: %s\n", loginName, request.Username)
}

func refresh(w http.ResponseWriter, r *http.Request) {
	claims, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRefreshHandlerError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddAuthToken(claims.Username, w.Header())
	log.Infof("%s: %s\n", refreshName, claims.Username)
}

func hashPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return nil, wrapHashPasswordError(err)
	}

	return hash, nil
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

func generateStarterPokemons(pokemonNr int) map[string]pokemons.Pokemon { // TODO only for testing

	toReturn := make(map[string]pokemons.Pokemon, pokemonNr)

	for i := 0; i < pokemonNr; i++ {
		newPokemon := pokemons.GetOneWildPokemon(
			5,
			500,
			500,
			25,
			300,
			fmt.Sprintf("starter-%d", i))
		toReturn[newPokemon.Id] = *newPokemon
	}
	return toReturn
}
