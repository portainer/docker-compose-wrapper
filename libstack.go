package libstack

import (
	"context"
)

type Deployer interface {
	Deploy(ctx context.Context, workingDir, host, projectName string, filePaths []string, envFilePath string, forceRecreate bool) error
	Remove(ctx context.Context, workingDir, host, projectName string, filePaths []string, envFilePath string) error
	Pull(ctx context.Context, workingDir, host, projectName string, filePaths []string, envFilePath string) error
}
