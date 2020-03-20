package auth

import (
	"github.com/NOVAPokemon/utils"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

// Types
type Claims struct {
	Id       primitive.ObjectID
	Username string
	Trainer  utils.Trainer
	jwt.StandardClaims
}

// Constants
const TokenCookieName = "token"
const JWTDuration = 30 * time.Minute

// Global variables
var JWTKey = []byte("my_secret_key")

func VerifyJWT(w *http.ResponseWriter, r *http.Request, caller string) (err error, claims *Claims) {
	c, err := r.Cookie(TokenCookieName)

	if err != nil {
		utils.HandleCookieError(w, caller, err)
		return err, nil
	}

	tknStr := c.Value
	claims = &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})

	if err != nil {
		utils.HandleJWTVerifyingError(w, caller, err)
		return err, nil
	}

	if !tkn.Valid && time.Unix(claims.ExpiresAt, 0).Unix() < time.Now().Unix() {
		(*w).WriteHeader(http.StatusUnauthorized)
		return err, nil
	}

	return nil, claims
}
