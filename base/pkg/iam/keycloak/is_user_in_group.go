package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
)

func (i *keyCloakIAM) IsUserInGroup(
	ctx context.Context,
	userID string,
	groupID string,
) (bool, error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return false, err
	}
	users, err := i.kc.GetGroupMembers(
		ctx,
		accessToken,
		i.kc.Realm,
		groupID,
		gocloak.GetGroupsParams{
			Max:    gocloak.IntP(1),
			Search: gocloak.StringP(userID),
		},
	)
	if err != nil {
		return false, errors.Wrap(err, "keycloak.GetGroupMembers")
	}
	for _, user := range users {
		if gocloak.PString(user.ID) == userID {
			return true, nil
		}
	}
	return false, nil
}
