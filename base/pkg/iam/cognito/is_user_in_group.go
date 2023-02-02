package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

func (i *cognitoIAM) IsUserInGroup(
	ctx context.Context,
	userID string,
	groupID string,
) (bool, error) {
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
			return false, errors.Wrap(err, "cognito.ListUsersInGroup")
		}
		for _, user := range result.Users {
			if aws.StringValue(user.Username) == userID {
				return true, nil
			}
		}
		nextToken = result.NextToken
		if nextToken == nil {
			return false, nil
		}
	}
}
