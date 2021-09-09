package composebinary

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	libstack "github.com/portainer/docker-compose-wrapper"
	liberrors "github.com/portainer/docker-compose-wrapper/compose/errors"
	"github.com/portainer/docker-compose-wrapper/compose/internal/utils"
)

// ComposeWrapper provide a type for managing docker compose commands
type ComposeWrapper struct {
	binaryPath string
	configPath string
}

// NewComposeWrapper initializes a new ComposeWrapper service with local docker-compose binary.
func NewComposeWrapper(binaryPath, configPath string) (libstack.Deployer, error) {
	if !utils.IsBinaryPresent(utils.ProgramPath(binaryPath, "docker-compose")) {
		return nil, liberrors.ErrBinaryNotFound
	}

	return &ComposeWrapper{binaryPath: binaryPath, configPath: configPath}, nil
}

// Deploy creates and starts containers
func (wrapper *ComposeWrapper) Deploy(ctx context.Context, workingDir, host, projectName string, filePaths []string, envFilePath string) error {
	output, err := wrapper.Command(newUpCommand(filePaths), workingDir, host, projectName, envFilePath)
	if len(output) != 0 {
		log.Printf("[libstack,composebinary] [message: finish deploying] [output: %s] [err: %s]", output, err)
	}

	return err
}

// Remove stops and removes containers
func (wrapper *ComposeWrapper) Remove(ctx context.Context, workingDir, host, projectName string, filePaths []string) error {
	output, err := wrapper.Command(newDownCommand(filePaths), workingDir, host, projectName, "")
	if len(output) != 0 {
		log.Printf("[libstack,composebinary] [message: finish deploying] [output: %s] [err: %s]", output, err)
	}

	return err
}

// Command executes a docker-compose command
func (wrapper *ComposeWrapper) Command(command composeCommand, workingDir, host, projectName, envFilePath string) ([]byte, error) {
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
	cmd := exec.Command(program, command.ToArgs()...)
	cmd.Dir = workingDir

	if wrapper.configPath != "" {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_CONFIG=%s", wrapper.configPath))
	}

	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, stderr.String())
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

func (command *composeCommand) WithHost(host string) {
	command.args = append(command.args, "-H", host)
}

func (command *composeCommand) ToArgs() []string {
	return append(command.args, command.command...)
}
