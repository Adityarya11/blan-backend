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

func RunSource(sourceCode string) (string, error) {
	fileName := fmt.Sprintf("%s.bl", uuid.New().String())
	filePath := filepath.Join("workspace", fileName)

	err := os.WriteFile(filePath, []byte(sourceCode), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write the codefile: %v", err)
	}

	defer os.Remove(filePath)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	binaryName := "blan"
	if runtime.GOOS == "windows" {
		binaryName = "blan.exe"
	}
	binaryPath := filepath.Join(".", binaryName)

	cmd := exec.CommandContext(ctx, binaryPath, filePath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out

	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("execution timeout: possible infinite loop")
	}

	if err != nil {
		return "", fmt.Errorf("%s", stderr.String())
	}

	return out.String(), nil

}
