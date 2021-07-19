package composebinary

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	libstack "github.com/portainer/docker-compose-wrapper"
)

var (
	// ErrBinaryNotFound is returned when docker-compose binary is not found
	ErrBinaryNotFound = errors.New("docker-compose binary not found")
)

// ComposeWrapper provide a type for managing docker compose commands
type ComposeWrapper struct {
	binaryPath string
	configPath string
}

// NewComposeWrapper initializes a new ComposeWrapper service with local docker-compose binary.
func NewComposeWrapper(binaryPath, configPath string) (libstack.Deployer, error) {
	if !IsBinaryPresent(programPath(binaryPath, "docker-compose")) {
		return nil, ErrBinaryNotFound
	}

	return &ComposeWrapper{binaryPath: binaryPath, configPath: configPath}, nil
}

// Up create and start containers
func (wrapper *ComposeWrapper) Deploy(ctx context.Context, host, projectName string, filePaths []string, envFilePath string) error {
	output, err := wrapper.Command(newUpCommand(filePaths), host, projectName, envFilePath, wrapper.configPath)
	if len(output) != 0 {
		log.Printf("[libstack,composebinary] [message: finish deploying] [output: %s] [err: %s]", output, err)
	}

	return err
}

// Down stop and remove containers
func (wrapper *ComposeWrapper) Remove(ctx context.Context, host, projectName string, filePaths []string) error {
	output, err := wrapper.Command(newDownCommand(filePaths), host, projectName, "", "")
	if len(output) != 0 {
		log.Printf("[libstack,composebinary] [message: finish deploying] [output: %s] [err: %s]", output, err)
	}

	return err
}

// Command exectue a docker-compose commanåd
func (wrapper *ComposeWrapper) Command(command composeCommand, url, projectName, envFilePath, configPath string) ([]byte, error) {
	program := programPath(wrapper.binaryPath, "docker-compose")

	if projectName != "" {
		command.WithProjectName(projectName)
	}

	if envFilePath != "" {
		command.WithEnvFilePath(envFilePath)
	}

	if url != "" {
		command.WithURL(url)
	}

	var stderr bytes.Buffer
	cmd := exec.Command(program, command.ToArgs()...)

	if configPath != "" {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_CONFIG=%s", configPath))
	}

	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	return output, nil
}

type composeCommand struct {
	command []string
	args    []string
}

func newCommand(command []string, filePaths []string) composeCommand {
	var args []string
	for _, path := range filePaths {
		args = append(args, "-f")
		args = append(args, strings.TrimSpace(path))
	}
	return composeCommand{
		args:    args,
		command: command,
	}
}

func newUpCommand(filePaths []string) composeCommand {
	return newCommand([]string{"up", "-d"}, filePaths)
}

func newDownCommand(filePaths []string) composeCommand {
	return newCommand([]string{"down", "--remove-orphans"}, filePaths)
}

func (command *composeCommand) WithProjectName(projectName string) {
	command.args = append(command.args, "-p", projectName)
}

func (command *composeCommand) WithEnvFilePath(envFilePath string) {
	command.args = append(command.args, "--env-file", envFilePath)
}

func (command *composeCommand) WithURL(url string) {
	command.args = append(command.args, "-H", url)
}

func (command *composeCommand) ToArgs() []string {
	return append(command.args, command.command...)
}
