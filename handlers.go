package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/clients"
	"github.com/NOVAPokemon/utils/cookies"
	userdb "github.com/NOVAPokemon/utils/database/user"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Constants
const StatusOnline = "online"

var trainersClient *clients.TrainersClient

func init() {
	trainersClient = clients.NewTrainersClient(fmt.Sprintf("%s:%d", host, utils.TrainersPort), &cookiejar.Jar{})
}

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
		Items:    generateStarterItems(6),
		Pokemons: generateStarterPokemons(6),
		Stats: utils.TrainerStats{
			Level: 0,
			Coins: 0,
		},
	}

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

	err, user := userdb.GetUserByUsername(request.Username)

	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	if !verifyPassword([]byte(request.Password), user.PasswordHash) {
		utils.HandleWrongPasswordError(&w, LoginName, request.Username, request.Password)
		return
	}

	err = cookies.SetAuthToken(request.Username, LoginName, &w)

	if err != nil {
		log.Error(err)
		http.Error(w, "Error occurred setting cookie", http.StatusInternalServerError)
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

func generateStarterItems(itemsNr int) map[string]utils.Item { //TODO only for testing

	toReturn := make(map[string]utils.Item, itemsNr)

	for i := 0; i < itemsNr; i++ {
		newItem := utils.Item{
			Id:   primitive.NewObjectID(),
			Name: fmt.Sprintf("item-%d", i),
		}
		toReturn[newItem.Id.Hex()] = newItem
	}

	return toReturn

}

func generateStarterPokemons(pokemonNr int) map[string]utils.Pokemon { //TODO only for testing

	toReturn := make(map[string]utils.Pokemon, pokemonNr)

	for i := 0; i < pokemonNr; i++ {
		newPokemon := utils.Pokemon{
			Id:      primitive.NewObjectID(),
			Species: fmt.Sprintf("species-%d", i),
			Damage:  10,
			Level:   0,
		}
		toReturn[newPokemon.Id.Hex()] = newPokemon
	}

	return toReturn

}
