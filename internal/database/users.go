package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"sort"
	"time"
)

var userCounter int = 1

type DatabaseUser struct {
	Id                    int       `json:"id"`
	Email                 string    `json:"email"`
	Password              []byte    `json:"password"`
	IsChirpyRed           bool      `json:"is_chirpy_red"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshExpirationDays time.Time `json:"refresh_expiration_days"`
}

type RespondUser struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (db *DB) CreateUser(email string, password []byte) (RespondUser, error) {
	if email == "" {
		return RespondUser{}, errors.New("cant create a User with no email")
	}

	NewUser := DatabaseUser{
		Id:       userCounter,
		Email:    email,
		Password: password,
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

func (db *DB) StoreRTokenAndExpiration(rTokenString string, userId int) error {
	storage, err := db.LoadDB()
	if err != nil {
		log.Println(err)
		return err
	}
	refreshExpiration := time.Now().UTC().Add(time.Hour * 1440)

	specificUser := storage.Users[userId]
	specificUser.RefreshToken = rTokenString
	specificUser.RefreshExpirationDays = refreshExpiration
	storage.Users[userId] = specificUser
	writeErr := db.writeDB(storage)
	if writeErr != nil {
		log.Println(writeErr)
		return writeErr
	}
	return nil
}

func (db *DB) RevokeRToken(rTokenString string) error {
	allUsers, err := db.GetUsers()
	if err != nil {
		log.Println(err)
		return err
	}
	var specificUser DatabaseUser
	var match bool
	for _, user := range allUsers {
		if user.RefreshToken == rTokenString {
			specificUser = user
			match = true
		}
	}
	if match {
		storage, err := db.LoadDB()
		if err != nil {
			log.Println(err)
			return err
		}
		specificUser.RefreshToken = ""
		storage.Users[specificUser.Id] = specificUser
		writeErr := db.writeDB(storage)
		if writeErr != nil {
			log.Println(writeErr)
			return writeErr
		}
		return nil

	} else {
		return errors.New("could not find a User with this specific RefreshToken")
	}

}

func (db *DB) AssignMemberSatusByID(userId int) error {
	storage, err := db.LoadDB()
	if err != nil {
		log.Println(err)
		return err
	}
	specificUser, ok := storage.Users[userId]
	if !ok {
		return errors.New("user does not exist")
	}
	specificUser.IsChirpyRed = true
	storage.Users[userId] = specificUser
	writeErr := db.writeDB(storage)

	if writeErr != nil {
		log.Println(writeErr)
		return writeErr
	}
	return nil
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
