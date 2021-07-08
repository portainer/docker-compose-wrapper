package libstack

import (
	"context"
)

type Deployer interface {
	Deploy(ctx context.Context, host, projectName string, filePaths []string, envFilePath string) error
	Remove(ctx context.Context, host, projectName string, filePaths []string) error
}
