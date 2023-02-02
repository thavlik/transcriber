package keycloak

import (
	"context"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) SetPassword(
	username string,
	password string,
	temporary bool,
) error {
	accessToken, err := i.refreshAccessToken(context.Background())
	if err != nil {
		return err
	}
	userID, err := i.resolveUser(context.Background(), username)
	if err != nil {
		return err
	}
	err = i.kc.SetPassword(
		context.Background(),
		accessToken,
		userID,
		i.kc.Realm,
		password,
		temporary,
	)
	if kcErr, ok := err.(*gocloak.APIError); ok && kcErr.Code == http.StatusNotFound {
		return iam.ErrUserNotFound
	} else if err != nil {
		return errors.Wrap(err, "keycloak")
	}
	return nil
}
