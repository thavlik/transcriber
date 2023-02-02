package base

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type CognitoOptions struct {
	UserPoolID               string
	ClientID                 string
	ClientSecret             string
	AllowTokenUseBeforeIssue bool
}

func (o *CognitoOptions) IsSet() bool {
	return o.UserPoolID != "" && o.ClientID != "" && o.ClientSecret != ""
}

func (o *CognitoOptions) Ensure() {
	if o.UserPoolID == "" {
		panic(errors.New("missing --cognito-user-pool-id"))
	}
	if o.ClientID == "" {
		panic(errors.New("missing --cognito-client-id"))
	}
	if o.ClientSecret == "" {
		panic(errors.New("missing --cognito-client-secret"))
	}
}

func CognitoEnv(o *CognitoOptions, required bool) *CognitoOptions {
	if o == nil {
		o = &CognitoOptions{}
	}
	CheckEnv("COGNITO_USER_POOL_ID", &o.UserPoolID)
	CheckEnv("COGNITO_CLIENT_ID", &o.ClientID)
	CheckEnv("COGNITO_CLIENT_SECRET", &o.ClientSecret)
	CheckEnvBool("COGNITO_ALLOW_TOKEN_USE_BEFORE_ISSUE", &o.AllowTokenUseBeforeIssue)
	if required {
		o.Ensure()
	}
	return o
}

func AddCognitoFlags(cmd *cobra.Command, o *CognitoOptions) {
	cmd.PersistentFlags().StringVar(&o.UserPoolID, "cognito-user-pool-id", "", "cognito user pool id")
	cmd.PersistentFlags().StringVar(&o.ClientID, "cognito-client-id", "", "cognito client id")
	cmd.PersistentFlags().StringVar(&o.ClientSecret, "cognito-client-secret", "", "cognito client secret")
	cmd.PersistentFlags().BoolVar(&o.AllowTokenUseBeforeIssue, "cognito-allow-token-use-before-issue", false, "permit use of cognito tokens with timestamps occuring before they are issued (bug workaround)")
}
