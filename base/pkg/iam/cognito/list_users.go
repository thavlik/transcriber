package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *cognitoIAM) ListUsers(ctx context.Context) ([]*iam.User, error) {
	var users []*iam.User
	var paginationToken *string
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		result, err := i.cognito.ListUsers(&cognitoidentityprovider.ListUsersInput{
			UserPoolId:      aws.String(i.userPoolID),
			PaginationToken: paginationToken,
		})
		if err != nil {
			return nil, errors.Wrap(err, "cognito")
		}
		for _, user := range result.Users {
			u := &iam.User{
				ID:       aws.StringValue(user.Username),
				Username: aws.StringValue(user.Username),
			}
			applyUserAttributes(u, user.Attributes)
			users = append(users, u)
		}
		paginationToken = result.PaginationToken
		if paginationToken == nil {
			break
		}
	}
	return users, nil
}
