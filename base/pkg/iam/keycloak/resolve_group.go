package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) ResolveGroup(
	ctx context.Context,
	groupName string,
) (string, error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return "", err
	}
	groups, err := i.kc.GetGroups(
		context.Background(),
		accessToken,
		i.kc.Realm,
		gocloak.GetGroupsParams{
			Max:    gocloak.IntP(1),
			Search: gocloak.StringP(groupName),
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "keycloak.GetGroups")
	}
	for _, group := range groups {
		if gocloak.PString(group.Name) == groupName {
			return gocloak.PString(group.ID), nil
		}
	}
	return "", iam.ErrGroupNotFound
}
