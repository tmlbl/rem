// Stores the last known state of the instance for each environment
package storage

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

// DB is a local data store for keeping track of last known instance state
type DB struct {
	rootPath string
}

func New(rootPath string) (*DB, error) {
	return &DB{
		rootPath: rootPath,
	}, nil
}

func Default() *DB {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(home, ".rem")
	err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	db, err := New(path)
	if err != nil {
		panic(err)
	}
	return db
}

type State struct {
	InstanceID  string
	BaseImageID string
	IP          net.IP
}

// SaveState stores the instance state information in the local database
func (db *DB) SaveState(path string, state *State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	path = filepath.Join(db.rootPath, path)
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

// GetState returns the last known state of the instance for the config file at
// the given path
func (db *DB) GetState(path string) (*State, error) {
	data, err := ioutil.ReadFile(filepath.Join(db.rootPath, path))
	if err != nil {
		return nil, err
	}
	state := State{}
	err = json.Unmarshal(data, &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}
