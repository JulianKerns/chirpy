package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (cfg *apiConfig) createToken(expiresAt time.Duration, userId int) (string, error) {
	expiringTime := expiresAt
	if expiresAt == time.Second*0 || expiresAt > time.Second*3600 {
		expiringTime = time.Second * 3600
	}
	userIdString := strconv.Itoa(userId)

	now := time.Now().UTC()
	expires := now.Add(expiringTime)

	if now.IsZero() || expires.IsZero() {
		return "", errors.New("time components are failing")
	}
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expires),
		Subject:   userIdString,
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		RegisteredClaims: claims,
	})

	if newToken == nil {
		return "", errors.New("could not create a Token")
	}
	signedToken, err := newToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (cfg *apiConfig) verifyToken(stringToken string) (*jwt.Token, error) {
	var userClaim UserClaims
	verifiedToken, err := jwt.ParseWithClaims(stringToken, &userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		//validator := jwt.NewValidator()
		//validateError := validator.Validate(verifiedToken.Claims)
		//if validateError != nil {
		//	if errors.Is(validateError, jwt.ErrTokenExpired) {
		//		return nil, errors.New("token is expired")
		//	} else if errors.Is(validateError, jwt.ErrTokenUnverifiable) {
		//		return nil, errors.New("bad token")
		//	}
		//	return nil, err
		//}
		return nil, err
	}
	if !verifiedToken.Valid {
		log.Println("invalid Token")
		return nil, err
	}
	return verifiedToken, nil
}

func (cfg *apiConfig) ValidateTokenGetId(w http.ResponseWriter, sentToken string) (int, error) {
	strippedToken, ok := strings.CutPrefix(sentToken, "Bearer ")
	if !ok {
		log.Println("could not remove prefix")
	}

	verifiedToken, err := cfg.verifyToken(strippedToken)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusUnauthorized, "Bad Authentication Token")
		return 0, err
	}
	userIdString, err := verifiedToken.Claims.GetSubject()
	if err != nil {
		log.Println("could not get the user id")
		log.Println(err)
		return 0, err
	}
	userIdInt, err := strconv.Atoi(userIdString)
	if err != nil {
		log.Println("could not transform the id strng to an int")
		log.Println(err)
		return 0, err
	}
	return userIdInt, nil
}
