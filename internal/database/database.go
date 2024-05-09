package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBstructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

var idCounter int = 1

func increase(i *int) int {
	*i++
	return *i

}

func NewDB(path string) (*DB, error) {
	_, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			os.WriteFile(path, []byte(""), 0100644)
		}
	}
	mutex := &sync.RWMutex{}
	newDB := &DB{
		path: path,
		mux:  mutex,
	}
	return newDB, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	if body == "" {
		return Chirp{}, errors.New("cant create a Chirp with an empty body")
	}
	NewChirp := Chirp{
		Id:   idCounter,
		Body: body,
	}
	content, err := json.Marshal(NewChirp)
	if err != nil {
		return Chirp{}, err
	}
	os.WriteFile(db.path, content, 0100644)
	storage, err := db.loadDB()

	if err != nil {
		log.Fatalf("%s", err)
		return Chirp{}, err
	}

	storage.Chirps[idCounter] = NewChirp
	fmt.Println(storage)
	increase(&idCounter)
	db.writeDB(storage)

	return NewChirp, nil

}

func (db *DB) GetChirps() ([]Chirp, error) {
	storage, err := db.loadDB()
	allChirps := []Chirp{}
	if err != nil {
		return []Chirp{}, err
	}
	for _, v := range storage.Chirps {
		allChirps = append(allChirps, v)
	}
	return allChirps, nil
}

//func (db *DB) ensure() error {
//	_, err := os.ReadFile(db.path)
//	if err != nil {
//		if errors.Is(err, os.ErrNotExist) {
//			os.WriteFile(db.path, []byte(""), 0100644)
//			return nil
//		}
//	}
//	return err
//}

func (db *DB) writeDB(dbStructure DBstructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	os.WriteFile(db.path, data, 0100644)
	return nil
}

func (db *DB) loadDB() (DBstructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := os.ReadFile(db.path)
	var databaseEntry = make(map[int]Chirp)
	database := DBstructure{databaseEntry}
	if data == nil {
		return database, nil
	}

	if err != nil {
		return DBstructure{}, err
	}
	errDB := json.Unmarshal(data, &database)
	if errDB != nil {
		return DBstructure{}, errDB
	}
	return database, nil
}
