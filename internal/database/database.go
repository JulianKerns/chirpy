package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	id   int
	body string
}

type DBstructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

var idCounter int = 1

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
	newChirp := Chirp{
		id:   idCounter,
		body: body,
	}
	storage, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	storage.Chirps[idCounter] = newChirp
	idCounter++

	db.writeDB(storage)

	return newChirp, nil

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
	db.mux.RLock()
	defer db.mux.Unlock()
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	os.WriteFile(db.path, data, 0100644)
	return nil
}

func (db *DB) loadDB() (DBstructure, error) {
	db.mux.RLock()
	defer db.mux.Unlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBstructure{}, err
	}
	database := DBstructure{}
	errDB := json.Unmarshal(data, &database)
	if errDB != nil {
		return DBstructure{}, errDB
	}
	return database, nil
}
