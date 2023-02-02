package keycloak

import (
	"context"

	"github.com/pkg/errors"
)

func (i *keyCloakIAM) AddUserToGroup(
	userID string,
	groupID string,
) error {
	accessToken, err := i.refreshAccessToken(context.Background())
	if err != nil {
		return err
	}
	if err := i.kc.AddUserToGroup(
		context.Background(),
		accessToken,
		i.kc.Realm,
		userID,
		groupID,
	); err != nil {
		return errors.Wrap(err, "keycloak")
	}
	return nil
}
