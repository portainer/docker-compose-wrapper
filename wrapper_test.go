package wrapper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setup(t *testing.T) *ComposeWrapper {
	w, err := NewComposeWrapper("")
	if err != nil {
		t.Fatal(err)
	}

	return w
}

func Test_UpAndDown(t *testing.T) {

	const composeFileContent = `version: "3.9"
services:
  busybox:
    image: "alpine:latest"
    container_name: "compose_wrapper_test"`
	const composedContainerName = "compose_wrapper_test"

	w := setup(t)

	dir := os.TempDir()

	filePath, err := createComposeFile(dir, composeFileContent)
	if err != nil {
		t.Fatal(err)
	}

	_, err := w.Up(filePath, "", "test1", "")
	if err != nil {
		t.Fatal(err)
	}

	if !containerExists(composedContainerName) {
		t.Fatal("container should exist")
	}

	_, err = w.Down(filePath, "", "test1")
	if err != nil {
		t.Fatal(err)
	}

	if containerExists(composedContainerName) {
		t.Fatal("container should be removed")
	}
}

type composeOptions struct {
	filePath    string
	url         string
	envFile     string
	projectName string
}

func createFile(dir, fileName, content string) (string, error) {
	filePath := filepath.Join(dir, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	f.WriteString(content)
	f.Close()

	return filePath, nil
}

func createEnvFile(dir, envFileContent string) (string, error) {
	return createFile(dir, "stack.env", envFileContent)
}

func createComposeFile(dir, composeFileContent string) (string, error) {
	return createFile(dir, "docmer-compose.yml", composeFileContent)
}

func containerExists(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-a", "-f", fmt.Sprintf("name=%s", containerName))

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to list containers: %s", err)
	}

	return strings.Contains(string(out), containerName)
}
