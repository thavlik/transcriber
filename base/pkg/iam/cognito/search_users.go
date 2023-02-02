package cognito

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *cognitoIAM) SearchUsers(
	ctx context.Context,
	prefix string,
) ([]*iam.User, error) {
	filter := aws.String(fmt.Sprintf(
		"username ^= \"%s\"",
		prefix,
	))
	var users []*iam.User
	var paginationToken *string
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		result, err := i.cognito.ListUsers(&cognitoidentityprovider.ListUsersInput{
			UserPoolId:      aws.String(i.userPoolID),
			PaginationToken: paginationToken,
			Limit:           aws.Int64(1),
			Filter:          filter,
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
