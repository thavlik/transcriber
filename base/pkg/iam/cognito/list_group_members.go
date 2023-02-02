package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *cognitoIAM) ListGroupMembers(
	ctx context.Context,
	groupID string,
) (users []*iam.User, err error) {
	var nextToken *string
	for {
		result, err := i.cognito.ListUsersInGroup(
			&cognitoidentityprovider.ListUsersInGroupInput{
				UserPoolId: aws.String(i.userPoolID),
				GroupName:  aws.String(groupID),
				NextToken:  nextToken,
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "cognito.ListUsersInGroup")
		}
		for _, user := range result.Users {
			v := &iam.User{
				ID:       aws.StringValue(user.Username),
				Username: aws.StringValue(user.Username),
			}
			applyUserAttributes(v, user.Attributes)
			users = append(users, v)
		}
		nextToken = result.NextToken
		if nextToken == nil {
			return users, nil
		}
	}
}
