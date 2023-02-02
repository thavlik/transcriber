package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

func (i *cognitoIAM) RemoveUserFromGroup(
	userID string,
	groupID string,
) error {
	if _, err := i.cognito.AdminRemoveUserFromGroup(
		&cognitoidentityprovider.AdminRemoveUserFromGroupInput{
			UserPoolId: aws.String(i.userPoolID),
			GroupName:  aws.String(groupID),
			Username:   aws.String(userID),
		},
	); err != nil {
		return errors.Wrap(err, "cognito")
	}
	return nil
}
