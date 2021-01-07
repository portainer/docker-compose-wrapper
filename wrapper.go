package wrapper

import (
	"bytes"
	"errors"
	"os/exec"
)

// ComposeWrapper provide a type for managing docker compose commands
type ComposeWrapper struct {
	binaryPath string
}

// NewComposeWrapper initializes a new ComposeWrapper service with local docker-compose binary.
func NewComposeWrapper() *ComposeWrapper {

	return &ComposeWrapper{binaryPath: ""}
}

// Up create and start containers
func (wrapper *ComposeWrapper) Up(filePath string) ([]byte, error) {
	return wrapper.Command([]string{"up", "-d"}, filePath)
}

// Down stop and remove containers
func (wrapper *ComposeWrapper) Down(filePath string) ([]byte, error) {
	return wrapper.Command([]string{"down", "--remove-orphans"}, filePath)
}

// Command exectue a docker-compose commanåd
func (wrapper *ComposeWrapper) Command(args []string, filePath string) ([]byte, error) {
	program := programPath(wrapper.binaryPath, "docker-compose")

	args = append([]string{"-f", filePath}, args...)

	var stderr bytes.Buffer
	cmd := exec.Command(program, args...)
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	return output, nil
}
