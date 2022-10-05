// Stores the last known state of the instance for each environment
package storage

import "fmt"

type State struct {
	InstanceID  string
	BaseImageID string
}

// GetState returns the last known state of the instance for the config file at
// the given path
func GetState(path string) (*State, error) {
	return nil, fmt.Errorf("not implemented")
}
