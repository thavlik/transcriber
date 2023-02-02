package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *cognitoIAM) SetPassword(
	username string,
	password string,
	temporary bool,
) error {
	if _, err := i.cognito.AdminSetUserPassword(
		&cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(i.userPoolID),
			Username:   aws.String(username),
			Password:   aws.String(password),
			Permanent:  aws.Bool(!temporary),
		},
	); err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == cognitoidentityprovider.ErrCodeUserNotFoundException {
			return iam.ErrUserNotFound
		}
		return errors.Wrap(err, "cognito")
	}
	return nil
}
