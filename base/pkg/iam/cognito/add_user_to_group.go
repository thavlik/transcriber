package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

func (i *cognitoIAM) AddUserToGroup(
	userID string,
	groupID string,
) error {
	if _, err := i.cognito.AdminAddUserToGroup(
		&cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(i.userPoolID),
			GroupName:  aws.String(groupID),
			Username:   aws.String(userID),
		},
	); err != nil {
		return errors.Wrap(err, "cognito")
	}
	return nil
}
