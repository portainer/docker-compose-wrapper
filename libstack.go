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
}

type DeployOptions struct {
	Options
	ForceRecreate bool
}
