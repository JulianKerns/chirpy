package main

import (
	"encoding/json"
	"log"
	"net/http"

	database "github.com/JulianKerns/chirpy/internal/database"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
		respondWithJSON(w, http.StatusOK, database.RespondUser{
			Id:    specificUser.Id,
			Email: specificUser.Email,
		})

	} else {
		respondError(w, http.StatusUnauthorized, "no User under this email present, or wrong Password")
	}
}