package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

// INPUT_URL=https://www.youtube.com/@weshammer runuser -pptruser -- node -e "$(cat /scripts/query-channel.js)"

func nodeQuery(
	ctx context.Context,
	scriptPath string,
	inputURL string,
	dest interface{},
) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	command := fmt.Sprintf(
		`INPUT_URL="%s" runuser -u pptruser -- node -e "$(cat %s)"`,
		inputURL,
		scriptPath,
	)
	cmd := exec.Command("bash", "-c", command)
	var stdout bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	select {
	case <-ctx.Done():
		if err := cmd.Process.Kill(); err != nil {
			panic(errors.Wrap(err, "failed to kill puppeteer child process"))
		}
		return errors.Errorf("context error: %v, hint: this could be because a waitForSelector operation timed out", ctx.Err())
	case err := <-done:
		if err != nil {
			return errors.Wrap(err, "run")
		}
	}
	if err := json.Unmarshal(stdout.Bytes(), &dest); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	return nil
}
