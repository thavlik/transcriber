package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *cognitoIAM) CreateGroup(name string) (*iam.Group, error) {
	if _, err := i.cognito.CreateGroup(&cognitoidentityprovider.CreateGroupInput{
		GroupName:  aws.String(name),
		UserPoolId: aws.String(i.userPoolID),
	}); err != nil {
		return nil, errors.Wrap(err, "cognito")
	}
	return &iam.Group{
		ID:   name,
		Name: name,
	}, nil
}
