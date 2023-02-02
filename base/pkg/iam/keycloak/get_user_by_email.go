package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) GetUserByEmail(
	ctx context.Context,
	email string,
) (*iam.User, error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	users, err := i.kc.GetUsers(
		ctx,
		accessToken,
		i.kc.Realm,
		gocloak.GetUsersParams{
			Email: gocloak.StringP(email),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "keycloak")
	}
	for _, user := range users {
		if gocloak.PString(user.Email) == email {
			return convertUser(user), nil
		}
	}
	return nil, iam.ErrUserNotFound
}
