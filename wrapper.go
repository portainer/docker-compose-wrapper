package wrapper

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	// ErrBinaryNotFound is returned when docker-compose binary is not found
	ErrBinaryNotFound = errors.New("docker-compose binary not found")
)

// ComposeWrapper provide a type for managing docker compose commands
type ComposeWrapper struct {
	binaryPath       string
	useComposePlugin bool
}

// NewComposeWrapper initializes a new ComposeWrapper service with local docker-compose binary.
func NewComposeWrapper(binaryPath string) (*ComposeWrapper, error) {
	dockerBinary := programPath(binaryPath, "docker")

	cliPluginsPath := path.Join(binaryPath, ".docker", "cli-plugins")
	composePlugin := programPath(cliPluginsPath, "docker-compose")

	usePlugins := IsBinaryPresent(dockerBinary) && IsBinaryPresent(composePlugin)
	if !usePlugins {
		program := programPath(binaryPath, "docker-compose")

		if !IsBinaryPresent(program) {
			return nil, ErrBinaryNotFound
		}
	}

	return &ComposeWrapper{binaryPath: binaryPath, useComposePlugin: usePlugins}, nil
}

// Up create and start containers
func (wrapper *ComposeWrapper) Up(filePaths []string, url, projectName, envFilePath, configPath string) ([]byte, error) {
	return wrapper.Command(newUpCommand(filePaths), url, projectName, envFilePath, configPath)
}

// Down stop and remove containers
func (wrapper *ComposeWrapper) Down(filePaths []string, url, projectName string) ([]byte, error) {
	return wrapper.Command(newDownCommand(filePaths), url, projectName, "", "")
}

// Command exectue a docker-compose commanåd
func (wrapper *ComposeWrapper) Command(command composeCommand, url, projectName, envFilePath, configPath string) ([]byte, error) {
	if projectName != "" {
		command.WithProjectName(projectName)
	}

	if envFilePath != "" {
		command.WithEnvFilePath(envFilePath)
	}

	if url != "" {
		command.WithURL(url)
	}

	program := programPath(wrapper.binaryPath, "docker-compose")
	args := command.ToArgs()
	if wrapper.useComposePlugin {
		log.Print("[DEBUG] [docker-compose-wrapper] [message: running docker with compose cli plugin]")
		program = programPath(wrapper.binaryPath, "docker")
		args = append([]string{"--config", path.Join(wrapper.binaryPath, ".docker"), "compose"}, args...)
	}

	cmd := exec.Command(program, args...)

	if configPath != "" && !wrapper.useComposePlugin {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_CONFIG=%s", configPath))
	}

	var stderr bytes.Buffer
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
		args = append(args, "-f", strings.TrimSpace(path))
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
