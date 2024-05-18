package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"sort"
)

var userCounter int = 1

type DatabaseUser struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Password     []byte `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type RespondUser struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refersh_token"`
}

func (db *DB) CreateUser(email string, password []byte) (RespondUser, error) {
	if email == "" {
		return RespondUser{}, errors.New("cant create a User with no email")
	}

	refreshToken, err := db.GenerateRefreshTokenString()
	if err != nil {
		return RespondUser{}, err
	}

	NewUser := DatabaseUser{
		Id:           userCounter,
		Email:        email,
		Password:     password,
		RefreshToken: refreshToken,
	}
	RUser := RespondUser{
		Id:    userCounter,
		Email: email,
	}
	storage, err := db.LoadDB()
	if err != nil {
		log.Printf("%s", err)
		return RespondUser{}, err
	}
	existingUsers, err := db.GetUsers()
	if err != nil {
		return RespondUser{}, err
	}

	for _, users := range existingUsers {
		if users.Email == email {
			return RespondUser{}, errors.New("email already exists")
		}

	}

	storage.Users[userCounter] = NewUser
	increase(&userCounter)
	db.writeDB(storage)
	return RUser, nil
}

func (db *DB) GenerateRefreshTokenString() (string, error) {
	sliceOfBytes := make([]byte, 32)
	_, err := rand.Read(sliceOfBytes)
	if err != nil {
		log.Println("could not read the random string")
		return "", err
	}
	byteToString := hex.EncodeToString(sliceOfBytes)
	return byteToString, nil
}

func (db *DB) GetUsers() ([]DatabaseUser, error) {
	storage, err := db.LoadDB()
	if err != nil {
		return []DatabaseUser{}, err
	}
	allUsers := []DatabaseUser{}
	for _, v := range storage.Users {
		allUsers = append(allUsers, v)
	}
	sort.Slice(allUsers, func(i, j int) bool { return allUsers[i].Id < allUsers[j].Id })

	return allUsers, nil
}

func (db *DB) GetUsersbyID(id int) (DatabaseUser, error) {
	storage, err := db.LoadDB()
	if err != nil {
		return DatabaseUser{}, err
	}
	var specificUser DatabaseUser
	var match bool
	for _, v := range storage.Users {
		if v.Id == id {
			specificUser = v
			match = true
		}
	}
	if match {
		return specificUser, nil
	} else {
		return DatabaseUser{}, errors.New("user does not exist")
	}

}

func (db *DB) UpdateUserByID(id int, email string, password []byte) error {
	storage, err := db.LoadDB()
	if err != nil {
		return err
	}
	user, err := db.GetUsersbyID(id)
	if err != nil {
		log.Println("user does not exist")
		return err
	}
	user.Email = email
	user.Password = password
	storage.Users[user.Id] = user
	db.writeDB(storage)
	return nil
}
