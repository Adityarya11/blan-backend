package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	cmd := exec.CommandContext(ctx, "./blan.exe", filePath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out

	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("CHUDDI!, Execution timeout, infinite loop hai re lodu.")
	}

	if err != nil {
		return "", fmt.Errorf("%s", stderr.String())
	}

	return out.String(), nil

}
