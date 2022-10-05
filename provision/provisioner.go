package provision

import (
	"fmt"

	"github.com/tmlbl/rem/config"
	"github.com/tmlbl/rem/provision/aws"
)

// Provisioner defines the behavior of a backend for a given provisioning
// service
type Provisioner interface {
	Build(base *config.Base) error
}

func GetProvisioner(name config.Platform) (Provisioner, error) {
	switch name {
	case config.PlatformAWS:
		return &aws.Provisioner{}, nil
	}
	return nil, fmt.Errorf("no provisioner for platform %s", name)
}
