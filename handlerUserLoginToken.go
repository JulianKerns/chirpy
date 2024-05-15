package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	database "github.com/JulianKerns/chirpy/internal/database"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) loginUserToken(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string        `json:"email"`
		Password         string        `json:"password"`
		ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Couldnt decode the Parameters")
		return
	}
	allUsers, err := cfg.DB.GetUsers()
	if err != nil {
		log.Println("could not retrieve the Users from the database")
	}
	var specificUser database.DatabaseUser
	var match bool
	for _, user := range allUsers {
		errCompare := bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
		if errCompare == nil {
			specificUser = user
			match = true
		}

	}
	if match {
		jwtToken, err := cfg.createToken(params.ExpiresInSeconds, specificUser.Id)
		if err != nil {
			log.Println(err)
			return
		}
		signedToken, err := cfg.signToken(jwtToken)
		if err != nil {
			log.Println(err)
			return
		}
		respondWithJSON(w, http.StatusOK, database.RespondUser{
			Id:    specificUser.Id,
			Email: specificUser.Email,
			Token: signedToken,
		})

	} else {
		respondError(w, http.StatusUnauthorized, "no User under this email present, or wrong Password")
	}
}
