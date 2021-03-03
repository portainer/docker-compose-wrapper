package libcomposewrapper

import (
	"context"

	"github.com/portainer/libcompose/config"
	"github.com/portainer/libcompose/docker"
	"github.com/portainer/libcompose/docker/ctx"
	"github.com/portainer/libcompose/lookup"
	"github.com/portainer/libcompose/project"
	"github.com/portainer/libcompose/project/options"
)

const (
	dockerClientVersion     = "1.24"
	composeSyntaxMaxVersion = "2"
)

// Service represents a service for managing compose stacks.
type Service struct {
}

// NewService initializes a new ComposeStackManager service.
func NewService() *Service {
	return &Service{}
}

// func (manager *ComposeStackManager) createClient(endpoint *portainer.Endpoint) (client.Factory, error) {

// 	endpointURL := endpoint.URL
// 	if endpoint.Type == portainer.EdgeAgentOnDockerEnvironment {
// 		tunnel := manager.reverseTunnelService.GetTunnelDetails(endpoint.ID)
// 		endpointURL = fmt.Sprintf("tcp://127.0.0.1:%d", tunnel.Port)
// 	}

// 	clientOpts := client.Options{
// 		Host:       endpointURL,
// 		APIVersion: dockerClientVersion,
// 	}

// 	if endpoint.TLSConfig.TLS {
// 		clientOpts.TLS = endpoint.TLSConfig.TLS
// 		clientOpts.TLSVerify = !endpoint.TLSConfig.TLSSkipVerify
// 		clientOpts.TLSCAFile = endpoint.TLSConfig.TLSCACertPath
// 		clientOpts.TLSCertFile = endpoint.TLSConfig.TLSCertPath
// 		clientOpts.TLSKeyFile = endpoint.TLSConfig.TLSKeyPath
// 	}

// 	return client.NewDefaultFactory(clientOpts)
// }

// ComposeSyntaxMaxVersion returns the maximum supported version of the libcompose syntax
func (manager *Service) ComposeSyntaxMaxVersion() string {
	return composeSyntaxMaxVersion
}

// Up will deploy a compose stack (equivalent of docker-compose up)
func (manager *Service) Up(filePath, url, projectName, envFilePath string) ([]byte, error) {

	// clientFactory, err := manager.createClient(endpoint)
	// if err != nil {
	// 	return err
	// }

	// env := make(map[string]string)
	// for _, envvar := range stack.Env {
	// 	env[envvar.Name] = envvar.Value
	// }

	proj, err := docker.NewProject(&ctx.Context{

		Context: project.Context{
			ComposeFiles: []string{filePath},
			EnvironmentLookup: &lookup.ComposableEnvLookup{
				Lookups: []config.EnvironmentLookup{
					&lookup.EnvfileLookup{
						Path: envFilePath,
					},
				},
			},
			ProjectName: projectName,
		},
	}, nil)

	if err != nil {
		return []byte(""), err
	}

	return []byte(""), proj.Up(context.Background(), options.Up{})
}

// Down will shutdown a compose stack (equivalent of docker-compose down)
func (manager *Service) Down(filePath, url, projectName string) ([]byte, error) {
	// clientFactory, err := manager.createClient(endpoint)
	// if err != nil {
	// 	return err
	// }

	proj, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{filePath},
			ProjectName:  projectName,
		},
	}, nil)

	if err != nil {
		return []byte(""), err
	}

	return []byte(""), proj.Down(context.Background(), options.Down{RemoveVolume: false, RemoveOrphans: true})
}
