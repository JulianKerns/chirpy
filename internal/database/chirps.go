package database

import (
	"errors"
	"log"
	"sync"
)

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

var idCounter int = 1

func (db *DB) CreateChirp(body string) (Chirp, error) {
	if body == "" {
		return Chirp{}, errors.New("cant create a Chirp with an empty body")
	}
	NewChirp := Chirp{
		Id:   idCounter,
		Body: body,
	}

	storage, err := db.LoadDB()
	if err != nil {
		log.Printf("%s", err)
		return Chirp{}, err
	}

	storage.Chirps[idCounter] = NewChirp
	increase(&idCounter)
	db.writeDB(storage)

	return NewChirp, nil

}

func (db *DB) CreateUserChirp(body string, userId int) (Chirp, error) {
	if body == "" {
		return Chirp{}, errors.New("cant create a Chirp with an empty body")
	}
	NewChirp := Chirp{
		Id:       idCounter,
		Body:     body,
		AuthorId: userId,
	}

	storage, err := db.LoadDB()
	if err != nil {
		log.Printf("%s", err)
		return Chirp{}, err
	}

	storage.Chirps[idCounter] = NewChirp
	increase(&idCounter)
	db.writeDB(storage)

	return NewChirp, nil
}
func (db *DB) GetUserChirps(authorId int) ([]Chirp, error) {
	storage, err := db.LoadDB()
	allUserChirps := []Chirp{}
	var chirpFound bool
	if err != nil {
		return []Chirp{}, err
	}
	for _, chirp := range storage.Chirps {

		if chirp.AuthorId == authorId {
			allUserChirps = append(allUserChirps, chirp)
			chirpFound = true
		}
	}
	if chirpFound {
		//sort.Slice(allUserChirps, func(i, j int) bool { return allUserChirps[i].Id < allUserChirps[j].Id })
		return allUserChirps, nil

	} else {
		return allUserChirps, errors.New("no Chirps found for this user")
	}
}
func (db *DB) GetChirps() ([]Chirp, error) {
	storage, err := db.LoadDB()
	allChirps := []Chirp{}
	if err != nil {
		return []Chirp{}, err
	}
	for _, v := range storage.Chirps {
		allChirps = append(allChirps, v)
	}
	//sort.Slice(allChirps, func(i, j int) bool { return allChirps[i].Id < allChirps[j].Id })

	return allChirps, nil
}

func (db *DB) DeleteUserChirp(authorId, chirpId int) error {
	userChirps, err := db.GetUserChirps(authorId)
	if err != nil {
		log.Println(err)
		return err
	}

	var match bool
	for _, chirp := range userChirps {
		if chirp.Id == chirpId {
			match = true
		}
	}
	if match {
		storage, err := db.LoadDB()
		if err != nil {
			log.Printf("%s", err)
			return err
		}
		delete(storage.Chirps, chirpId)
		db.writeDB(storage)
		return nil

	} else {
		return errors.New("no authorization to delete this chirp")
	}
}
