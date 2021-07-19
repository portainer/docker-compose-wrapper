package composebinary

import (
	"os/exec"
	"path"
	"runtime"
)

func osProgram(program string) string {
	if runtime.GOOS == "windows" {
		program += ".exe"
	}
	return program
}

func programPath(rootPath, program string) string {
	return path.Join(rootPath, osProgram(program))
}

// IsBinaryPresent check if docker compose binary is present
func IsBinaryPresent(program string) bool {
	_, err := exec.LookPath(program)
	return err == nil
}
