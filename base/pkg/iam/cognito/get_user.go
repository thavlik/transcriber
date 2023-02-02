package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func applyUserAttributes(
	user *iam.User,
	attrs []*cognitoidentityprovider.AttributeType,
) {
	for _, attr := range attrs {
		switch aws.StringValue(attr.Name) {
		case "email":
			user.Email = aws.StringValue(attr.Value)
		case "given_name":
			user.FirstName = aws.StringValue(attr.Value)
		case "family_name":
			user.LastName = aws.StringValue(attr.Value)
		}
	}
}

func (i *cognitoIAM) GetUser(
	ctx context.Context,
	username string,
) (*iam.User, error) {
	user, err := i.cognito.AdminGetUser(&cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(i.userPoolID),
		Username:   aws.String(username),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == cognitoidentityprovider.ErrCodeUserNotFoundException {
			return nil, iam.ErrUserNotFound
		}
		return nil, errors.Wrap(err, "cognito")
	}
	u := &iam.User{
		ID:       username,
		Username: username,
		Enabled:  aws.BoolValue(user.Enabled),
	}
	applyUserAttributes(u, user.UserAttributes)
	return u, nil
}
