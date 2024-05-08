package database

import (
	"errors"
	"io/fs"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

//type Chirp struct {
//	id   int
//	body string
//}
//type DBstructure struct {
//	Chirps map[int]Chirp `json:"chirps"`
//}

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

//func (db *DB) ensure() error {
//	_, err := os.ReadFile(db.path)
//	if err == os.ErrNotExist {
//		os.WriteFile(db.path, []byte(""), 0100644)
//	}
//	return nil
//}
