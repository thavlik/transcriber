package cognito

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *cognitoIAM) DeleteGroup(groupID string) error {
	if _, err := i.cognito.DeleteGroup(&cognitoidentityprovider.DeleteGroupInput{
		GroupName:  aws.String(groupID),
		UserPoolId: aws.String(i.userPoolID),
	}); err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == cognitoidentityprovider.ErrCodeResourceNotFoundException {
			return iam.ErrGroupNotFound
		}
		return errors.Wrap(err, "cognito")
	}
	return nil
}
