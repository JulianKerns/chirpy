package database

import (
	"errors"
	"log"
	"sort"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
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

func (db *DB) GetChirps() ([]Chirp, error) {
	storage, err := db.LoadDB()
	allChirps := []Chirp{}
	if err != nil {
		return []Chirp{}, err
	}
	for _, v := range storage.Chirps {
		allChirps = append(allChirps, v)
	}
	sort.Slice(allChirps, func(i, j int) bool { return allChirps[i].Id < allChirps[j].Id })

	return allChirps, nil
}
