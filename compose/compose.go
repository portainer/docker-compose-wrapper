package compose

import (
	"context"
	"fmt"
	"log"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose-cli/pkg/api"
	"github.com/docker/compose-cli/pkg/compose"
	libstack "github.com/portainer/docker-compose-wrapper"
	"github.com/portainer/docker-compose-wrapper/compose/internal/composebinary"
)

type ComposeDeployer struct {
	configPath string
}

// NewComposeDeployer will try to create a wrapper for docker-compose binary
// if it's not availbale will use compose.Deployer
func NewComposeDeployer(binaryPath, configPath string) (libstack.Deployer, error) {
	deployer, err := composebinary.NewComposeWrapper(binaryPath, configPath)
	if err == nil {
		return deployer, nil
	}

	if err == composebinary.ErrBinaryNotFound {
		log.Printf("[INFO] [main,compose] [message: binary is missing, falling-back to compose library] [error: %s]", err)
		return &ComposeDeployer{configPath: configPath}, nil
	}

	return nil, err
}

// Up creates and starts containers
func (deployer *ComposeDeployer) Deploy(ctx context.Context, host, projectName string, filePaths []string, envFilePath string) error {
	service, err := prepareService(host, deployer.configPath)
	if err != nil {
		return fmt.Errorf("failed creating compose service: %w", err)
	}

	project, err := prepareProject(filePaths, projectName, envFilePath)
	if err != nil {
		return fmt.Errorf("failed preparing project: %w", err)
	}

	err = service.Up(ctx, project, api.UpOptions{})
	if err != nil {
		return fmt.Errorf("failed deploying: %w", err)
	}

	return nil
}

// Down stops and removes containers
func (deployer *ComposeDeployer) Remove(ctx context.Context, host, projectName string, filePaths []string) error {
	service, err := prepareService(host, deployer.configPath)
	if err != nil {
		return fmt.Errorf("failed creating compose service: %w", err)
	}

	err = service.Down(ctx, projectName, api.DownOptions{RemoveOrphans: true})
	if err != nil {
		return fmt.Errorf("failed removing: %w", err)
	}

	return nil
}

func prepareService(host, configPath string) (api.Service, error) {
	// compose-go outputs to stdout/stderr anyway, there's no way to overwrite it for now.
	dockercli, err := command.NewDockerCli(command.WithStandardStreams())
	if err != nil {
		return nil, fmt.Errorf("error creating client %w", err)
	}

	initOpts := flags.NewClientOptions()
	if host != "" {
		initOpts.Common.Hosts = []string{host}
	}
	if configPath != "" {
		initOpts.ConfigDir = configPath
	}

	err = dockercli.Initialize(initOpts)
	if err != nil {
		return nil, fmt.Errorf("error init client %w", err)
	}

	return compose.NewComposeService(dockercli.Client(), dockercli.ConfigFile()), nil
}

func prepareProject(filePaths []string, projectName, envFilePath string) (*types.Project, error) {
	additionalOpts := []cli.ProjectOptionsFn{cli.WithName(projectName)}

	if envFilePath != "" {
		additionalOpts = append(additionalOpts, cli.WithEnvFile(envFilePath))
	}

	projectOpts, err := cli.NewProjectOptions(filePaths, additionalOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed creating project options: %w", err)
	}

	project, err := cli.ProjectFromOptions(projectOpts)
	if err != nil {
		return nil, fmt.Errorf("error loading files: %w", err)
	}

	return project, nil
}
