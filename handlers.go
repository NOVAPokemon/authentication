package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/authentication/auth"
	"github.com/NOVAPokemon/utils"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
	userdb "github.com/NOVAPokemon/utils/database/user"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Constants
const StatusOnline = "online"

// Indicates if the server is online.
func Status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, StatusOnline)
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

	trainer, err := trainerdb.GetTrainerByUsername(request.Username)

	expirationTime := time.Now().Add(auth.JWTDuration)
	claims := &auth.Claims{
		Username:       request.Username,
		Trainer:        *trainer,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(auth.JWTKey)

	if err != nil {
		utils.HandleJWTSigningError(&w, LoginName, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    auth.TokenCookieName,
		Value:   tokenString,
		Expires: expirationTime,
	})

	log.Infof("%s: %s\n", LoginName, request.Username)
}

// Endpoint to refresh the token. Expects the user to already have a token.
func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(auth.TokenCookieName)

	if err != nil {
		utils.HandleCookieError(&w, RefreshName, err)
		return
	}

	tknStr := c.Value
	claims := &auth.Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.JWTKey, nil
	})

	if err != nil {
		utils.HandleJWTVerifyingError(&w, RefreshName, err)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 2*time.Minute {
		utils.HandleToSoonToRefreshError(&w, RefreshName)
		return
	}

	// Now, create a new token for the current user, with a renewed expiration time
	expirationTime := time.Now().Add(auth.JWTDuration)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(auth.JWTKey)

	if err != nil {
		utils.HandleJWTSigningError(&w, RefreshName, err)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    auth.TokenCookieName,
		Value:   tokenString,
		Expires: expirationTime,
	})

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
