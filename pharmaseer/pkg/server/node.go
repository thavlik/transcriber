package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// INPUT_URL=https://www.youtube.com/@weshammer runuser -pptruser -- node -e "$(cat /scripts/query-channel.js)"

func nodeQuery(scriptPath string, input string, dest interface{}) error {
	command := fmt.Sprintf(`INPUT_URL="%s" runuser -u pptruser -- node -e "$(cat %s)"`,
		input,
		scriptPath,
	)
	cmd := exec.Command("bash", "-c",
		command)
	var stdout bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(stdout.String())
		return errors.Wrap(err, "node")
		//return fmt.Errorf("failed to run '%s': %v", command, err)
	}
	if err := json.Unmarshal(stdout.Bytes(), &dest); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	return nil
}
