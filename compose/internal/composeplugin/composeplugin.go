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
func (wrapper *PluginWrapper) Deploy(ctx context.Context, filePaths []string, options libstack.DeployOptions) error {
	output, err := wrapper.command(newUpCommand(filePaths, upOptions{
		forceRecreate:        options.ForceRecreate,
		abortOnContainerExit: options.AbortOnContainerExit,
	}), options.Options,
	)
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
func (wrapper *PluginWrapper) Remove(ctx context.Context, filePaths []string, options libstack.Options) error {
	output, err := wrapper.command(newDownCommand(filePaths), options)
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
func (wrapper *PluginWrapper) Pull(ctx context.Context, filePaths []string, options libstack.Options) error {
	output, err := wrapper.command(newPullCommand(filePaths), options)
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
func (wrapper *PluginWrapper) command(command composeCommand, options libstack.Options) ([]byte, error) {
	program := utils.ProgramPath(wrapper.binaryPath, "docker-compose")

	if options.ProjectName != "" {
		command.WithProjectName(options.ProjectName)
	}

	if options.EnvFilePath != "" {
		command.WithEnvFilePath(options.EnvFilePath)
	}

	if options.Host != "" {
		command.WithHost(options.Host)
	}

	var stderr bytes.Buffer

	args := []string{}
	args = append(args, command.ToArgs()...)

	cmd := exec.Command(program, args...)
	cmd.Dir = options.WorkingDir

	if wrapper.configPath != "" {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "DOCKER_CONFIG="+wrapper.configPath)
	}

	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		errOutput := stderr.String()
		log.Printf("[INFO] [message: docker compose command failed] [output: %s] [error_output: %s] [error: %s]", output, errOutput, err)
		// stderr output outputs useless information such as "Removing network stack_default"
		return nil, errors.WithMessage(err, "docker-compose command failed")
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

type upOptions struct {
	forceRecreate        bool
	abortOnContainerExit bool ``
}

func newUpCommand(filePaths []string, options upOptions) composeCommand {
	args := []string{"up"}

	if options.abortOnContainerExit {
		args = append(args, "--abort-on-container-exit")
	} else { // detach by default, not working with --abort-on-container-exit
		args = append(args, "-d")
	}

	if options.forceRecreate {
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
