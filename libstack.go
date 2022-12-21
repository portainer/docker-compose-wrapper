package libstack

import (
	"context"
)

type Deployer interface {
	Deploy(ctx context.Context, filePaths []string, options DeployOptions) error
	Remove(ctx context.Context, filePaths []string, options Options) error
	Pull(ctx context.Context, filePaths []string, options Options) error
}

type Options struct {
	WorkingDir  string
	Host        string
	ProjectName string
	// EnvFilePath is the path to a .env file
	EnvFilePath string
	// Env is a list of environment variables to pass to the command, example: "FOO=bar"
	Env             []string
	PotentialErrors []string
}

type DeployOptions struct {
	Options
	ForceRecreate bool
	// AbortOnContainerExit will stop the deployment if a container exits.
	// This is useful when running a onetime task.
	//
	// When this is set, docker compose will output its logs to stdout
	AbortOnContainerExit bool ``
}
