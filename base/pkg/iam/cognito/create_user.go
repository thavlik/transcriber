package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

func (i *cognitoIAM) CreateUser(
	user *iam.User,
	password string,
	temporary bool,
) (id string, err error) {
	if user.ID != "" && user.ID != user.Username {
		return "", errors.New("cognito requires user id must be unset or the same as username")
	}
	user.ID = user.Username
	if _, err := i.cognito.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		Username:          aws.String(user.Username),
		UserPoolId:        aws.String(i.userPoolID),
		TemporaryPassword: aws.String(password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{{
			Name:  aws.String("email"),
			Value: aws.String(user.Email),
		}, {
			Name:  aws.String("given_name"),
			Value: aws.String(user.FirstName),
		}, {
			Name:  aws.String("family_name"),
			Value: aws.String(user.LastName),
		}},
	}); err != nil {
		return "", errors.Wrap(err, "cognito")
	}
	if !temporary {
		if err := i.SetPassword(
			user.Username,
			password,
			false,
		); err != nil {
			return "", errors.Wrap(err, "SetPassword")
		}
	}
	i.log.Debug("created user",
		zap.String("user.ID", user.Username),
		zap.String("user.Username", user.Username),
		zap.String("user.FirstName", user.FirstName),
		zap.String("user.LastName", user.LastName),
		zap.String("user.Email", user.Email))
	return user.Username, nil
}
