package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) GetUser(ctx context.Context, username string) (*iam.User, error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	users, err := i.kc.GetUsers(
		ctx,
		accessToken,
		i.kc.Realm,
		gocloak.GetUsersParams{
			Username: gocloak.StringP(username),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "keycloak")
	}
	for _, user := range users {
		if gocloak.PString(user.Username) == username {
			return convertUser(user), nil
		}
	}
	return nil, iam.ErrUserNotFound
}
