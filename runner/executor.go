package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
)

func InitWorkspace() error {
	return os.MkdirAll("workspace", 0755)
}

func RunSource(sourceCode string) (string, string, error) {
	fileName := fmt.Sprintf("%s.bl", uuid.New().String())
	filePath := filepath.Join("workspace", fileName)

	err := os.WriteFile(filePath, []byte(sourceCode), 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to write the codefile: %v", err)
	}
	defer os.Remove(filePath)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	compilerPath := os.Getenv("BLAN_BINARY_PATH")

	if compilerPath == "" {

		exePath, err := os.Executable()
		if err != nil {
			return "", "", fmt.Errorf("failed to resolve server executable path: %w", err)
		}

		binaryName := "blan"
		if runtime.GOOS == "windows" {
			binaryName = "blan.exe"
		}

		compilerPath = filepath.Join(filepath.Dir(exePath), binaryName)
	}

	cmd := exec.CommandContext(ctx, compilerPath, filePath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()

	outStr := out.String()
	errStr := stderr.String()

	if ctx.Err() == context.DeadlineExceeded {
		return outStr, errStr, fmt.Errorf("execution timeout: possibly infinite loop...")
	}

	if err != nil {
		return outStr, errStr, fmt.Errorf("execution failed: %w", err)
	}

	return outStr, errStr, nil
}
