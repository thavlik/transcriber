package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *cognitoIAM) GetUserByEmail(
	ctx context.Context,
	email string,
) (*iam.User, error) {
	resp, err := i.cognito.ListUsers(&cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(i.userPoolID),
		Filter:     aws.String(fmt.Sprintf(`email = "%s"`, email)),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == cognitoidentityprovider.ErrCodeUserNotFoundException {
			return nil, iam.ErrUserNotFound
		}
		return nil, errors.Wrap(err, "cognito")
	}
	for _, user := range resp.Users {
		u := &iam.User{
			ID:       aws.StringValue(user.Username),
			Username: aws.StringValue(user.Username),
			Enabled:  aws.BoolValue(user.Enabled),
		}
		applyUserAttributes(u, user.Attributes)
		return u, nil
	}
	return nil, iam.ErrUserNotFound
}
