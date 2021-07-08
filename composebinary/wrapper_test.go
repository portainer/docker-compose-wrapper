package composebinary

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	libstack "github.com/portainer/docker-compose-wrapper"
)

func setup(t *testing.T) libstack.Deployer {
	w, err := NewComposeWrapper("", "")
	if err != nil {
		t.Fatal(err)
	}

	return w
}

func Test_NewCommand_SingleFilePath(t *testing.T) {
	cmd := newCommand([]string{"up", "-d"}, []string{"docker-compose.yml"})
	expected := []string{"-f", "docker-compose.yml"}
	if !reflect.DeepEqual(cmd.args, expected) {
		t.Errorf("wrong output args, want: %v, got: %v", expected, cmd.args)
	}
}

func Test_NewCommand_MultiFilePaths(t *testing.T) {
	cmd := newCommand([]string{"up", "-d"}, []string{"docker-compose.yml", "docker-compose-override.yml"})
	expected := []string{"-f", "docker-compose.yml", "-f", "docker-compose-override.yml"}
	if !reflect.DeepEqual(cmd.args, expected) {
		t.Errorf("wrong output args, want: %v, got: %v", expected, cmd.args)
	}
}

func Test_NewCommand_MultiFilePaths_WithSpaces(t *testing.T) {
	cmd := newCommand([]string{"up", "-d"}, []string{" docker-compose.yml", "docker-compose-override.yml "})
	expected := []string{"-f", "docker-compose.yml", "-f", "docker-compose-override.yml"}
	if !reflect.DeepEqual(cmd.args, expected) {
		t.Errorf("wrong output args, want: %v, got: %v", expected, cmd.args)
	}
}

func Test_UpAndDown(t *testing.T) {

	const composeFileContent = `version: "3.9"
services:
  busybox:
    image: "alpine:3.7"
    container_name: "test_container_one"`

	const overrideComposeFileContent = `version: "3.9"
services:
  busybox:
    image: "alpine:latest"
    container_name: "test_container_two"`

	const composeContainerName = "test_container_two"

	w := setup(t)

	dir := os.TempDir()
	defer os.RemoveAll(dir)

	filePathOriginal, err := createFile(dir, "docker-compose.yml", composeFileContent)
	if err != nil {
		t.Fatal(err)
	}

	filePathOverride, err := createFile(dir, "docker-compose-override.yml", overrideComposeFileContent)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	err = w.Deploy(ctx, "", "test1", []string{filePathOriginal, filePathOverride}, "")
	if err != nil {
		t.Fatal(err)
	}

	if !containerExists(composeContainerName) {
		t.Fatal("container should exist")
	}

	err = w.Remove(ctx, "", "test1", []string{filePathOriginal, filePathOverride})
	if err != nil {
		t.Fatal(err)
	}

	if containerExists(composeContainerName) {
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
