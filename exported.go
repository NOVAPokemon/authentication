package main

import (
	"github.com/NOVAPokemon/utils"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func VerifyJWT(w *http.ResponseWriter, r *http.Request) (err error, username string) {
	c, err := r.Cookie(TokenCookieName)

	if err != nil {
		utils.HandleCookieError(w, RefreshName, err)
		return err, ""
	}

	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		utils.HandleJWTVerifyingError(w, RefreshName, err)
		return err, ""
	}

	if !tkn.Valid && time.Unix(claims.ExpiresAt, 0).Unix() < time.Now().Unix() {
		(*w).WriteHeader(http.StatusUnauthorized)
		return err, ""
	}

	return nil, claims.Username
}
