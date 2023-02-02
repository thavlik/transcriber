package iam

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/thavlik/transcriber/base/pkg/iam/api"

	"github.com/spf13/cobra"
)

var loginArgs struct {
	endpoint string
	timeout  time.Duration
	username string
	password string
	decode   bool
	decoder  string
	split    bool
}

var loginCmd = &cobra.Command{
	Use: "login",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if loginArgs.username == "" {
			return errors.New("missing --username")
		}
		if loginArgs.password == "" {
			return errors.New("missing --password")
		}
		if loginArgs.decode && loginArgs.decoder == "" {
			loginArgs.decoder = "std"
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.NewRemoteIAMClient(
			loginArgs.endpoint,
			api.NewRemoteIAMClientOptions().
				SetTimeout(loginArgs.timeout),
		).Login(
			context.Background(),
			api.LoginRequest{
				Username: loginArgs.username,
				Password: loginArgs.password,
			},
		)
		if err != nil {
			return err
		}
		if !loginArgs.decode && loginArgs.decoder == "" {
			if loginArgs.split {
				parts := strings.Split(resp.AccessToken, ".")
				for _, part := range parts {
					fmt.Println(part)
				}
			} else {
				fmt.Println(resp.AccessToken)
			}
			return nil
		}
		if !loginArgs.split {
			result, err := encoding().DecodeString(resp.AccessToken)
			fmt.Println(string(result))
			if err != nil {
				return err
			}
			return nil
		}
		parts := strings.Split(resp.AccessToken, ".")
		fmt.Printf("access token has %d parts\n", len(parts))
		offset := 0
		for i, part := range parts {
			result, err := encoding().DecodeString(part)
			func() {
				fmt.Printf("\npart %d (offset %d) ===============\n", i, offset)
				defer fmt.Println("=================================")
				if err != nil {
					fmt.Println(part)
					fmt.Printf("error decoding part %d as %s base64: %v\n", i, loginArgs.decoder, err)
					return
				}
				fmt.Println(string(result))
			}()
			offset += len(part) + 1
		}
		return nil
	},
}

func encoding() *base64.Encoding {
	switch loginArgs.decoder {
	case "std":
		return base64.StdEncoding
	case "url":
		return base64.URLEncoding
	case "rawstd":
		return base64.RawStdEncoding
	case "rawurl":
		return base64.RawURLEncoding
	default:
		panic(fmt.Errorf("unrecognized base64 encoding '%s'", loginArgs.decoder))
	}
}

func init() {
	loginCmd.PersistentFlags().StringVar(&loginArgs.endpoint, "endpoint", "http://localhost:8080", "admin service endpoint")
	loginCmd.PersistentFlags().DurationVar(&loginArgs.timeout, "timeout", defaultTimeout, "admin service timeout")
	loginCmd.PersistentFlags().StringVar(&loginArgs.username, "username", "", "username")
	loginCmd.PersistentFlags().StringVar(&loginArgs.password, "password", "", "password")
	loginCmd.PersistentFlags().BoolVarP(&loginArgs.split, "split", "s", false, "split parts")
	loginCmd.PersistentFlags().BoolVarP(&loginArgs.decode, "decode", "d", false, "decode base64")
	loginCmd.PersistentFlags().StringVar(&loginArgs.decoder, "decoder", "", "name of base64 decoder [ std | url | rawstd | rawurl ] (default is std)")
	iamCmd.AddCommand(loginCmd)
}
