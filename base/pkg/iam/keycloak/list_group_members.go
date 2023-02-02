package keycloak

import (
	"context"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) ListGroupMembers(
	ctx context.Context,
	groupID string,
) ([]*iam.User, error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	users, err := i.kc.GetGroupMembers(
		ctx,
		accessToken,
		i.kc.Realm,
		groupID,
		gocloak.GetGroupsParams{},
	)
	if kcErr, ok := err.(*gocloak.APIError); ok && kcErr.Code == http.StatusNotFound {
		return nil, iam.ErrGroupNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "keycloak")
	}
	return convertUsers(users), nil
}
