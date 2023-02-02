package cognito

import (
	"context"
)

func (i *cognitoIAM) ResolveGroup(
	ctx context.Context,
	groupName string,
) (string, error) {
	// with cognito, the name is the id
	return groupName, nil
}
