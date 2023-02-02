package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
)

func (i *keyCloakIAM) DeleteGroupByName(name string) error {
	accessToken, err := i.refreshAccessToken(context.Background())
	if err != nil {
		return err
	}
	groups, err := i.kc.GetGroups(
		context.Background(),
		accessToken,
		i.kc.Realm,
		gocloak.GetGroupsParams{
			Search: gocloak.StringP(name),
		},
	)
	if err != nil {
		return errors.Wrap(err, "keycloak")
	}
	for _, group := range groups {
		if gocloak.PString(group.Name) == name {
			if err := i.kc.DeleteGroup(
				context.Background(),
				accessToken,
				i.kc.Realm,
				gocloak.PString(group.ID),
			); err != nil {
				return errors.Wrap(err, "keycloak")
			}
		}
	}
	return nil
}
