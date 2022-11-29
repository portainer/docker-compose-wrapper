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
	EnvFilePath string
	Env         map[string]string
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
