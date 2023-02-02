package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) CreateGroup(name string) (*iam.Group, error) {
	accessToken, err := i.refreshAccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	groupID, err := i.kc.CreateGroup(
		context.Background(),
		accessToken,
		i.kc.Realm,
		gocloak.Group{
			Name: gocloak.StringP(name),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "keycloak")
	}
	return &iam.Group{
		ID:   groupID,
		Name: name,
	}, nil
}
