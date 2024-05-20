package main

import (
	"encoding/json"

	"log"
	"net/http"

	database "github.com/JulianKerns/chirpy/internal/database"

	bcrypt "golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
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
	sentToken := r.Header.Get("Authorization")

	validatedUserId, err := cfg.ValidateTokenGetId(w, sentToken)
	if err != nil {
		log.Println(err)
		return
	}

	hash, errHash := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if errHash != nil {
		log.Fatalln("could not hash the password correctly")
	}
	errUpdate := cfg.DB.UpdateUserByID(validatedUserId, params.Email, []byte(hash))
	if errUpdate != nil {
		log.Println("could not update the user")
		respondError(w, http.StatusUnauthorized, "access to this user is not allowed")
	}
	respondWithJSON(w, http.StatusOK, database.RespondUser{
		Id:    validatedUserId,
		Email: params.Email,
	})

}
