package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
)

var isDiseaseArgs struct {
	base.ServiceOptions
}

var isDiseaseCmd = &cobra.Command{
	Use: "is-disease",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.ServiceEnv("", &isDiseaseArgs.ServiceOptions)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.TrimSpace(strings.Join(args, " "))
		if input == "" {
			return errors.New("missing input")
		}
		req, err := http.NewRequestWithContext(
			cmd.Context(),
			http.MethodGet,
			isDiseaseArgs.Endpoint+"/disease?q="+url.QueryEscape(input),
			nil,
		)
		if err != nil {
			return err
		}
		resp, err := (&http.Client{
			Timeout: isDiseaseArgs.Timeout,
		}).Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, body)
		}
		_, err = io.Copy(cmd.OutOrStdout(), resp.Body)
		return err
	},
}

func init() {
	rootCmd.AddCommand(isDiseaseCmd)
	base.AddServiceFlags(isDiseaseCmd, "", &isDiseaseArgs.ServiceOptions, 10*time.Second)
}
