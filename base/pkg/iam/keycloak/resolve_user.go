package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) resolveUser(
	ctx context.Context,
	username string,
) (userID string, err error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return "", err
	}
	users, err := i.kc.GetUsers(
		ctx,
		accessToken,
		i.kc.Realm,
		gocloak.GetUsersParams{
			Max:      gocloak.IntP(1),
			Username: gocloak.StringP(username),
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "keycloak")
	}
	for _, user := range users {
		if gocloak.PString(user.Username) == username {
			return gocloak.PString(user.ID), nil
		}
	}
	return "", iam.ErrUserNotFound
}
