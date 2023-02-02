package keycloak

import (
	"context"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

func (i *keyCloakIAM) DeleteUser(username string) error {
	accessToken, err := i.refreshAccessToken(context.Background())
	if err != nil {
		return err
	}
	userID, err := i.resolveUser(context.Background(), username)
	if err != nil {
		return err
	}
	err = i.kc.DeleteUser(
		context.Background(),
		accessToken,
		i.kc.Realm,
		userID,
	)
	if kcErr, ok := err.(*gocloak.APIError); ok && kcErr.Code == http.StatusNotFound {
		return iam.ErrUserNotFound
	} else if err != nil {
		return errors.Wrap(err, "keycloak")
	}
	i.log.Debug("deleted user",
		zap.String("realm", i.kc.Realm),
		zap.String("username", username),
		zap.String("userID", userID))
	return nil
}
