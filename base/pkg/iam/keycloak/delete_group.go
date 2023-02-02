package keycloak

import (
	"context"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) DeleteGroup(groupID string) error {
	accessToken, err := i.refreshAccessToken(context.Background())
	if err != nil {
		return err
	}
	err = i.kc.DeleteGroup(
		context.Background(),
		accessToken,
		i.kc.Realm,
		groupID,
	)
	if kcErr, ok := err.(*gocloak.APIError); ok && kcErr.Code == http.StatusNotFound {
		return iam.ErrGroupNotFound
	} else if err != nil {
		return errors.Wrap(err, "keycloak")
	}
	return nil
}
