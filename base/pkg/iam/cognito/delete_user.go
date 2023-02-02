package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

func (i *cognitoIAM) DeleteUser(username string) error {
	if _, err := i.cognito.AdminDeleteUser(&cognitoidentityprovider.AdminDeleteUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(i.userPoolID),
	}); err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == cognitoidentityprovider.ErrCodeUserNotFoundException {
			return iam.ErrUserNotFound
		}
		return errors.Wrap(err, "cognito")
	}
	i.log.Debug("deleted user",
		zap.String("username", username))
	return nil
}
