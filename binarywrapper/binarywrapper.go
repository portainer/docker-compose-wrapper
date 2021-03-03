package binarywrapper

import (
	"bytes"
	"errors"
	"os/exec"

	"github.com/portainer/docker-compose-wrapper/internal/utils"
)

// ComposeWrapper provide a type for managing docker compose commands
type ComposeWrapper struct {
	binaryPath       string
	syntaxMaxVersion string
}

// NewComposeWrapper initializes a new ComposeWrapper service with local docker-compose binary.
func NewComposeWrapper(binaryPath, syntaxMaxVersion string) (*ComposeWrapper, error) {
	return &ComposeWrapper{
		binaryPath:       binaryPath,
		syntaxMaxVersion: syntaxMaxVersion,
	}, nil
}

// ComposeSyntaxMaxVersion returns the maximum supported version of the docker compose syntax
func (wrapper *ComposeWrapper) ComposeSyntaxMaxVersion() string {
	return wrapper.syntaxMaxVersion
}

// Up create and start containers
func (wrapper *ComposeWrapper) Up(filePath, url, projectName, envFilePath string) ([]byte, error) {
	return wrapper.Command(newUpCommand(filePath), url, projectName, envFilePath)
}

// Down stop and remove containers
func (wrapper *ComposeWrapper) Down(filePath, url, projectName string) ([]byte, error) {
	return wrapper.Command(newDownCommand(filePath), url, projectName, "")
}

// Command exectue a docker-compose commanåd
func (wrapper *ComposeWrapper) Command(command composeCommand, url, projectName, envFilePath string) ([]byte, error) {
	program := utils.ProgramPath(wrapper.binaryPath, "docker-compose")

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

func newCommand(command []string, filePath string) composeCommand {
	return composeCommand{
		args:    []string{"-f", filePath},
		command: command,
	}
}

func newUpCommand(filePath string) composeCommand {
	return newCommand([]string{"up", "-d"}, filePath)
}

func newDownCommand(filePath string) composeCommand {
	return newCommand([]string{"down", "--remove-orphans"}, filePath)
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
