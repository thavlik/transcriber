package keycloak

import (
	"context"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) ListUserGroups(
	ctx context.Context,
	userID string,
) ([]*iam.Group, error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	groups, err := i.kc.GetUserGroups(
		ctx,
		accessToken,
		i.kc.Realm,
		userID,
		gocloak.GetGroupsParams{},
	)
	if kcErr, ok := err.(*gocloak.APIError); ok && kcErr.Code == http.StatusNotFound {
		return nil, iam.ErrUserNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "keycloak")
	}
	n := len(groups)
	result := make([]*iam.Group, n)
	for i, group := range groups {
		result[i] = &iam.Group{
			ID:   gocloak.PString(group.ID),
			Name: gocloak.PString(group.Name),
		}
	}
	return result, nil
}
