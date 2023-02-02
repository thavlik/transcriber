package base

import (
	"fmt"

	"github.com/spf13/cobra"
)

type IAMDriver string

var (
	CognitoDriver  IAMDriver = "cognito"
	KeyCloakDriver IAMDriver = "keycloak"
)

type IAMOptions struct {
	Driver   IAMDriver
	KeyCloak KeyCloakOptions
	Cognito  CognitoOptions
}

func IAMEnv(o *IAMOptions, required bool) {
	KeyCloakEnv(&o.KeyCloak, false)
	CognitoEnv(&o.Cognito, false)
	if required {
		switch o.Driver {
		case "":
			panic("missing --iam-driver")
		case CognitoDriver:
			o.Cognito.Ensure()
		case KeyCloakDriver:
			o.KeyCloak.Ensure()
		default:
			panic(fmt.Errorf("unrecognized iam driver '%s'", o.Driver))
		}
	}
}

func AddIAMFlags(cmd *cobra.Command, o *IAMOptions) {
	AddKeyCloakFlags(cmd, &o.KeyCloak)
	AddCognitoFlags(cmd, &o.Cognito)
	cmd.PersistentFlags().StringVar((*string)(&o.Driver), "iam-driver", "", "iam driver [ cognito | keycloak ]")
}
