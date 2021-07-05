package compose

import (
	"log"

	libstack "github.com/portainer/docker-compose-wrapper"
	"github.com/portainer/docker-compose-wrapper/compose/errors"
	"github.com/portainer/docker-compose-wrapper/compose/internal/composebinary"
	"github.com/portainer/docker-compose-wrapper/compose/internal/composeplugin"
)

// NewComposeDeployer will try to create a wrapper for docker-compose binary
// if it's not availbale it will try to create a wrapper for docker-compose plugin
func NewComposeDeployer(binaryPath, configPath string) (libstack.Deployer, error) {
	deployer, err := composebinary.NewComposeWrapper(binaryPath, configPath)
	if err == nil {
		return deployer, nil
	}

	if err == errors.ErrBinaryNotFound {
		log.Printf("[INFO] [main,compose] [message: binary is missing, falling-back to compose plugin] [error: %s]", err)
		return composeplugin.NewPluginWrapper(binaryPath, configPath)
	}

	return nil, err
}
