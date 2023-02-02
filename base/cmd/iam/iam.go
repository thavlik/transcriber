package iam

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"
	cognito_iam "github.com/thavlik/transcriber/base/pkg/iam/cognito"
	keycloak_iam "github.com/thavlik/transcriber/base/pkg/iam/keycloak"
	"go.uber.org/zap"
)

var defaultTimeout = 10 * time.Second

var iamCmd = &cobra.Command{
	Use: "iam",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please choose a subcommand")
	},
}

func AddIAMSubCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(iamCmd)
}

func InitIAM(o *base.IAMOptions, log *zap.Logger) iam.IAM {
	switch o.Driver {
	case base.CognitoDriver:
		return cognito_iam.NewCognitoIAM(
			o.Cognito.AllowTokenUseBeforeIssue,
			log,
		)
	case base.KeyCloakDriver:
		return keycloak_iam.NewKeyCloakIAM(
			base.ConnectKeyCloak(&o.KeyCloak),
			log,
		)
	default:
		panic(fmt.Errorf("unrecognized iam driver '%s'", o.Driver))
	}
}
