package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) SearchUsers(
	ctx context.Context,
	prefix string,
) ([]*iam.User, error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	users, err := i.kc.GetUsers(
		ctx,
		accessToken,
		i.kc.Realm,
		gocloak.GetUsersParams{
			Search: gocloak.StringP(prefix),
			Max:    gocloak.IntP(1),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "keycloak")
	}
	return convertUsers(users), nil
}
