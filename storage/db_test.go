package storage

import (
	"io/ioutil"
	"net"
	"testing"
)

func getTestDB(t *testing.T) *DB {
	dir, err := ioutil.TempDir("", "rem-test-db-*")
	if err != nil {
		t.Error(err)
	}
	db, err := New(dir)
	if err != nil {
		t.Error(err)
	}
	return db
}

func TestSaveState(t *testing.T) {
	testPath := "/home/tim/projects/foo/rem.yaml"
	state := State{
		InstanceID:  "test-1",
		BaseImageID: "base-image",
		IP:          net.ParseIP("0.0.0.0"),
	}

	db := getTestDB(t)
	err := db.SaveState(testPath, &state)
	if err != nil {
		t.Error(err)
	}

	state2, err := db.GetState(testPath)
	if err != nil {
		t.Error(err)
	}

	if state2.InstanceID != state.InstanceID ||
		state2.BaseImageID != state.BaseImageID {
		t.Errorf("structs did not match")
	}
}
