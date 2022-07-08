package composeplugin

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	libstack "github.com/portainer/docker-compose-wrapper"
	"github.com/portainer/docker-compose-wrapper/compose/internal/utils"
)

var (
	MissingDockerComposePluginErr = errors.New("docker-compose plugin is missing from config path")
)

// PluginWrapper provide a type for managing docker compose commands
type PluginWrapper struct {
	binaryPath string
	configPath string
}

// NewPluginWrapper initializes a new ComposeWrapper service with local docker-compose binary.
func NewPluginWrapper(binaryPath, configPath string) (libstack.Deployer, error) {
	if !utils.IsBinaryPresent(utils.ProgramPath(binaryPath, "docker-compose")) {
		return nil, MissingDockerComposePluginErr
	}

	return &PluginWrapper{binaryPath: binaryPath, configPath: configPath}, nil
}

// Up create and start containers
func (wrapper *PluginWrapper) Deploy(ctx context.Context, workingDir, host, projectName string, filePaths []string, envFilePath string, forceRecreate bool) error {
	output, err := wrapper.command(newUpCommand(filePaths, forceRecreate), workingDir, host, projectName, envFilePath)
	if len(output) != 0 {
		if err != nil {
			return err
		}

		log.Printf("[INFO] [docker compose] [message: Stack deployment successful]")
		log.Printf("[DEBUG] [docker compose] [output: %s]", output)
	}

	return err
}

// Down stop and remove containers
func (wrapper *PluginWrapper) Remove(ctx context.Context, workingDir, host, projectName string, filePaths []string, envFilePath string) error {
	output, err := wrapper.command(newDownCommand(filePaths), workingDir, host, projectName, envFilePath)
	if len(output) != 0 {
		if err != nil {
			return err
		}

		log.Printf("[INFO] [docker compose] [message: Stack removal successful]")
		log.Printf("[DEBUG] [docker compose] [output: %s]", output)
	}

	return err
}

// Pull images
func (wrapper *PluginWrapper) Pull(ctx context.Context, workingDir, host, projectName string, filePaths []string, envFilePath string) error {
	output, err := wrapper.command(newPullCommand(filePaths), workingDir, host, projectName, envFilePath)
	if len(output) != 0 {
		if err != nil {
			return err
		}

		log.Printf("[INFO] [docker compose] [message: Stack pull successful]")
		log.Printf("[DEBUG] [docker compose] [output: %s]", output)
	}

	return err
}

// Command execute a docker-compose command
func (wrapper *PluginWrapper) command(command composeCommand, workingDir, host, projectName, envFilePath string) ([]byte, error) {
	program := utils.ProgramPath(wrapper.binaryPath, "docker-compose")

	if projectName != "" {
		command.WithProjectName(projectName)
	}

	if envFilePath != "" {
		command.WithEnvFilePath(envFilePath)
	}

	if host != "" {
		command.WithHost(host)
	}

	var stderr bytes.Buffer

	args := []string{}
	args = append(args, command.ToArgs()...)

	cmd := exec.Command(program, args...)
	cmd.Dir = workingDir

	if wrapper.configPath != "" {
		if wrapper.configPath != "" {
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "DOCKER_CONFIG="+wrapper.configPath)
		}
	}

	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	return output, nil
}

type composeCommand struct {
	globalArgs        []string // docker-compose global arguments: --host host -f file.yaml
	subCommandAndArgs []string // docker-compose subcommand:  up, down folllowed by subcommand arguments
}

func newCommand(command []string, filePaths []string) composeCommand {
	args := []string{}
	for _, path := range filePaths {
		args = append(args, "-f")
		args = append(args, strings.TrimSpace(path))
	}
	return composeCommand{
		globalArgs:        args,
		subCommandAndArgs: command,
	}
}

func newUpCommand(filePaths []string, forceRereate bool) composeCommand {
	args := []string{"up", "-d"}
	if forceRereate {
		args = append(args, "--force-recreate")
	}
	return newCommand(args, filePaths)
}

func newDownCommand(filePaths []string) composeCommand {
	return newCommand([]string{"down", "--remove-orphans"}, filePaths)
}

func newPullCommand(filePaths []string) composeCommand {
	return newCommand([]string{"pull"}, filePaths)
}

func (command *composeCommand) WithHost(host string) {
	// prepend compatibility flags such as this one as they must appear before the
	// regular global args otherwise docker-compose will throw an error
	command.globalArgs = append([]string{"--host", host}, command.globalArgs...)
}

func (command *composeCommand) WithProjectName(projectName string) {
	command.globalArgs = append(command.globalArgs, "--project-name", projectName)
}

func (command *composeCommand) WithEnvFilePath(envFilePath string) {
	command.globalArgs = append(command.globalArgs, "--env-file", envFilePath)
}

func (command *composeCommand) ToArgs() []string {
	return append(command.globalArgs, command.subCommandAndArgs...)
}
