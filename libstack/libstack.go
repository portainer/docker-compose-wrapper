package libstack

type Deployer interface {
	Deploy(projectName, host string, filePaths []string, envFilePath string) error
	Remove(projectName, host string, filePaths []string) error
}
