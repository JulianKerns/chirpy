package main

import (
	"encoding/json"
	"log"
	"net/http"

	bcrypt "golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

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

	hash, errHash := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if errHash != nil {
		log.Fatalln("could not hash the password correctly")
	}

	newUser, errCh := cfg.DB.CreateUser(params.Email, hash)
	if errCh != nil {
		respondError(w, http.StatusConflict, "Email already exists")
		log.Println("could not add the User to the database, email already exists")

	} else {
		respondWithJSON(w, http.StatusCreated, newUser)
	}
}
