package keycloak

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/Nerzal/gocloak/v12"

	"github.com/pkg/errors"
)

func retrieveAuthHeader(refreshToken string) (*gocloak.JWT, error) {
	refreshToken = strings.ReplaceAll(refreshToken, "Bearer ", "")
	body, err := base64.URLEncoding.DecodeString(refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "base64")
	}
	jwt := &gocloak.JWT{}
	if err := json.Unmarshal(body, jwt); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return jwt, nil
}

func (i *keyCloakIAM) Logout(
	ctx context.Context,
	token string,
) error {
	jwt, err := retrieveAuthHeader(token)
	if err != nil {
		return err
	}
	if err := i.kc.Logout(
		ctx,
		i.kc.ClientID,
		i.kc.ClientSecret,
		i.kc.Realm,
		jwt.RefreshToken,
	); err != nil {
		return errors.Wrap(err, "keycloak")
	}
	return nil
}
