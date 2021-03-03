package wrapper

import (
	"github.com/portainer/docker-compose-wrapper/binarywrapper"
	"github.com/portainer/docker-compose-wrapper/internal/utils"
	"github.com/portainer/docker-compose-wrapper/libcomposewrapper"
)

// DockerComposeWrapper represents a service to manage Compose stacks
type DockerComposeWrapper interface {
	ComposeSyntaxMaxVersion() string
	Up(filePath, url, projectName, envFilePath string) ([]byte, error)
	Down(filePath, url, projectName string) ([]byte, error)
}

// NewDockerComposeWrapper initializes a new DockerComposeWrapper service with local docker-compose binary.
func NewDockerComposeWrapper(binaryPath, defaultSyntaxMaxVersion string) (DockerComposeWrapper, error) {
	if !utils.IsBinaryPresent(utils.ProgramPath(binaryPath, "docker-compose")) {
		return libcomposewrapper.NewService(), nil
	}

	return binarywrapper.NewComposeWrapper(binaryPath, defaultSyntaxMaxVersion)
}
