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

func TestCommand(t *testing.T) {
	w := setup(t)

	file := "docker-compose-test.yml"
	_, err := w.Up(file, "", "", "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Down(file, "", "")
	if err != nil {
		t.Fatal(err)
	}

}

const composeFile = `version: "3.9"
services:
  busybox:
    image: "alpine:latest"
    container_name: "compose_wrapper_test"`
const composedContainerName = "compose_wrapper_test"

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

// func Test_UpAndDown(t *testing.T) {

// 	stack, endpoint := setup(t)

// 	w, err := NewComposeStackManager("", nil)
// 	if err != nil {
// 		t.Fatalf("Failed creating manager: %s", err)
// 	}

// 	err = w.Up(stack, endpoint)
// 	if err != nil {
// 		t.Fatalf("Error calling docker-compose up: %s", err)
// 	}

// 	if !containerExists(composedContainerName) {
// 		t.Fatal("container should exist")
// 	}

// 	err = w.Down(stack, endpoint)
// 	if err != nil {
// 		t.Fatalf("Error calling docker-compose down: %s", err)
// 	}

// 	if containerExists(composedContainerName) {
// 		t.Fatal("container should be removed")
// 	}
// }

func containerExists(containerName string) bool {
	cmd := exec.Command("docker", "ps", "-a", "-f", fmt.Sprintf("name=%s", containerName))

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to list containers: %s", err)
	}

	return strings.Contains(string(out), containerName)
}
