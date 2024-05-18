package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"
)

type DBstructure struct {
	Chirps map[int]Chirp        `json:"chirps"`
	Users  map[int]DatabaseUser `json:"users"`
}

func increase(i *int) int {
	*i++
	return *i

}

func NewDB(path string) (*DB, error) {
	mutex := &sync.RWMutex{}
	newDB := &DB{
		path: path,
		mux:  mutex,
	}
	errDB := newDB.ensure()
	if errDB != nil {
		return &DB{}, errDB
	}
	return newDB, nil
}

func (db *DB) ensure() error {
	_, err := os.ReadFile(db.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			os.WriteFile(db.path, []byte("{}"), 0100644)
			return nil
		}
	}
	return err
}

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

func (db *DB) LoadDB() (DBstructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := os.ReadFile(db.path)
	var databaseEntry = make(map[int]Chirp)
	var databaseUsers = make(map[int]DatabaseUser)
	database := DBstructure{
		Chirps: databaseEntry,
		Users:  databaseUsers,
	}
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
