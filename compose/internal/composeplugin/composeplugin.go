package composeplugin

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	libstack "github.com/portainer/docker-compose-wrapper"
	liberrors "github.com/portainer/docker-compose-wrapper/compose/errors"
	"github.com/portainer/docker-compose-wrapper/compose/internal/utils"
)

var (
	MissingDockerComposePluginErr = errors.New("docker compose plugin is missing from config path")
)

// PluginWrapper provide a type for managing docker compose commands
type PluginWrapper struct {
	binaryPath string
	configPath string
}

// NewPluginWrapper initializes a new ComposeWrapper service with local docker binary.
func NewPluginWrapper(binaryPath, configPath string) (libstack.Deployer, error) {
	if !utils.IsBinaryPresent(utils.ProgramPath(binaryPath, "docker")) {
		return nil, liberrors.ErrBinaryNotFound
	}

	if configPath == "" {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.WithMessage(err, "failed fetching user home directory")
		}
		configPath = path.Join(homePath, ".docker")
	}

	dockerPluginsPath := path.Join(configPath, "cli-plugins")
	pluginPath := utils.ProgramPath(binaryPath, "docker-compose.plugin")

	if utils.IsBinaryPresent(pluginPath) {
		if !utils.IsBinaryPresent(utils.ProgramPath(dockerPluginsPath, "docker-compose")) {
			err := os.MkdirAll(dockerPluginsPath, 0755)
			if err != nil {
				return nil, errors.WithMessage(err, "failed creating plugins path")
			}
		}

		err := utils.Move(pluginPath, utils.ProgramPath(dockerPluginsPath, "docker-compose"))
		if err != nil {
			return nil, err
		}
	} else if !utils.IsBinaryPresent(utils.ProgramPath(dockerPluginsPath, "docker-compose")) {
		return nil, MissingDockerComposePluginErr
	}

	return &PluginWrapper{binaryPath: binaryPath, configPath: configPath}, nil
}

// Up create and start containers
func (wrapper *PluginWrapper) Deploy(ctx context.Context, workingDir, host, projectName string, filePaths []string, envFilePath string, forceRereate bool) error {
	output, err := wrapper.command(newUpCommand(filePaths, forceRereate), workingDir, host, projectName, envFilePath)
	if len(output) != 0 {
		log.Printf("[libstack,composebinary] [message: finish deploying] [output: %s] [err: %s]", output, err)
	}

	return err
}

// Down stop and remove containers
func (wrapper *PluginWrapper) Remove(ctx context.Context, workingDir, host, projectName string, filePaths []string) error {
	output, err := wrapper.command(newDownCommand(filePaths), workingDir, host, projectName, "")
	if len(output) != 0 {
		log.Printf("[libstack,composebinary] [message: finish deploying] [output: %s] [err: %s]", output, err)
	}

	return err
}

// Pull images
func (wrapper *PluginWrapper) Pull(ctx context.Context, workingDir, host, projectName string, filePaths []string) error {
	output, err := wrapper.command(newPullCommand(filePaths), workingDir, host, projectName, "")
	if len(output) != 0 {
		log.Printf("[libstack,composebinary] [message: finish pulling] [output: %s] [err: %s]", output, err)
	}

	return err
}

// Command exectue a docker-compose commanåd
func (wrapper *PluginWrapper) command(command composeCommand, workingDir, url, projectName, envFilePath string) ([]byte, error) {
	program := utils.ProgramPath(wrapper.binaryPath, "docker")

	if projectName != "" {
		command.WithProjectName(projectName)
	}

	if envFilePath != "" {
		command.WithEnvFilePath(envFilePath)
	}

	var stderr bytes.Buffer

	args := []string{}

	if url != "" {
		args = append(args, "-H", url)
	}

	args = append(args, "compose")
	args = append(args, command.ToArgs()...)

	cmd := exec.Command(program, args...)
	cmd.Dir = workingDir

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

//TO-DO: double check docker compose plugin logic
//to see if it is necessary to implement forceRereate arg here
func newUpCommand(filePaths []string, forceRereate bool) composeCommand {
	return newCommand([]string{"up", "-d"}, filePaths)
}

func newDownCommand(filePaths []string) composeCommand {
	return newCommand([]string{"down", "--remove-orphans"}, filePaths)
}

func newPullCommand(filePaths []string) composeCommand {
	return newCommand([]string{"pull"}, filePaths)
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
