package main

import (
	"errors"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
}

func (cfg *apiConfig) createToken(expiresAt time.Duration, userId int) (*jwt.Token, error) {
	expiringTime := expiresAt
	if expiresAt == 0 || expiresAt > 24 {
		expiringTime = 24
	}
	userIdString := strconv.Itoa(userId)
	claims := UserClaims{
		jwt.RegisteredClaims{
			Issuer:    "chripy",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiringTime)),
			Subject:   userIdString,
		},
	}
	var emptyToken *jwt.Token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if newToken == emptyToken {
		return emptyToken, errors.New("could not create a Token")
	}

	return newToken, nil
}

func (cfg *apiConfig) signToken(token *jwt.Token) (string, error) {
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (cfg *apiConfig) verifyToken(stringToken string) (*jwt.Token, error) {
	verifiedToken, err := jwt.ParseWithClaims(stringToken, UserClaims{}, jwt.Keyfunc(func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	}))
	var emptyToken *jwt.Token
	if err != nil {
		return emptyToken, err
	}
	return verifiedToken, nil
}
