package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/database/user"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Types
type Claims struct {
	Username string
	jwt.StandardClaims
}

// Constants
const StatusOnline = "online"
const TokenCookieName = "token"
const JWTDuration = 30 * time.Minute

// Global variables
var jwtKey = []byte("my_secret_key")

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

	userToAdd := utils.User{
		Id:        primitive.NewObjectID(),
		TrainerId: primitive.NewObjectID(),
		Username:  request.Username,
		PasswordHash: hashPassword([]byte(request.Password)),
	}

	err, id := user.AddUser(&userToAdd)

	if err != nil {
		return
	}

	log.Infof("%s: %s %s\n", RegisterName, request.Username, id)
}

// Logs in a user. Expects a JSON with username and password in the body.
func Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.HandleJSONDecodeError(&w, LoginName, err)
		return
	}

	// TEMPORARY while MongoDB API is not done
	err, user := user.GetUserByUsername(request.Username)

	if err != nil {
		return
	}

	if !verifyPassword([]byte(request.Password), user.PasswordHash) {
		utils.HandleWrongPasswordError(&w, LoginName, request.Username, request.Password)
		return
	}

	expirationTime := time.Now().Add(JWTDuration)
	claims := &Claims{
		request.Username,
		jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		utils.HandleJWTSigningError(&w, LoginName, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    TokenCookieName,
		Value:   tokenString,
		Expires: expirationTime,
	})

	log.Infof("%s: %s\n", LoginName, request.Username)
}

// Endpoint to refresh the token. Expects the user to already have a token.
func Refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(TokenCookieName)

	if err != nil {
		utils.HandleCookieError(&w, RefreshName, err)
		return
	}

	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
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
	expirationTime := time.Now().Add(JWTDuration)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		utils.HandleJWTSigningError(&w, RefreshName, err)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    TokenCookieName,
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