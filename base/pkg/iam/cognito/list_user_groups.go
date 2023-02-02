package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *cognitoIAM) ListUserGroups(
	ctx context.Context,
	userID string,
) ([]*iam.Group, error) {
	result, err := i.cognito.AdminListGroupsForUser(
		&cognitoidentityprovider.AdminListGroupsForUserInput{
			UserPoolId: aws.String(i.userPoolID),
			Username:   aws.String(userID),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "cognito.AdminListGroupsForUser")
	}
	n := len(result.Groups)
	groups := make([]*iam.Group, n)
	for i, group := range result.Groups {
		groups[i] = &iam.Group{
			ID:   aws.StringValue(group.GroupName),
			Name: aws.StringValue(group.GroupName),
		}
	}
	return groups, nil
}
