package cognito

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

type cognitoIAM struct {
	cognito                  *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID               string
	clientID                 string
	clientSecret             string
	region                   string
	allowTokenUseBeforeIssue bool
	log                      *zap.Logger
}

func NewCognitoIAM(
	allowUseBeforeIssue bool,
	log *zap.Logger,
) iam.IAM {
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		panic("missing AWS_REGION")
	}
	userPoolID, ok := os.LookupEnv("COGNITO_USER_POOL_ID")
	if !ok {
		panic("missing COGNITO_USER_POOL_ID")
	}
	clientID, ok := os.LookupEnv("COGNITO_CLIENT_ID")
	if !ok {
		panic("missing COGNITO_CLIENT_ID")
	}
	// some application may not have a client secret
	clientSecret := os.Getenv("COGNITO_CLIENT_SECRET")
	cognito := cognitoidentityprovider.New(
		session.Must(session.NewSession(
			&aws.Config{
				Region: aws.String(region),
			})))
	return &cognitoIAM{
		cognito,
		userPoolID,
		clientID,
		clientSecret,
		region,
		allowUseBeforeIssue,
		log,
	}
}
