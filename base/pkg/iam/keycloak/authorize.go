package keycloak

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

func (i *keyCloakIAM) Authorize(
	ctx context.Context,
	token string,
	permissions []string,
) (string, error) {
	jwt, err := retrieveAuthHeader(token)
	if err != nil {
		return "", err
	}
	result, err := i.kc.RetrospectToken(
		ctx,
		jwt.AccessToken,
		i.kc.ClientID,
		i.kc.ClientSecret,
		i.kc.Realm,
	)
	if err != nil {
		return "", errors.Wrap(err, "keycloak")
	} else if !gocloak.PBool(result.Active) {
		return "", iam.ErrTokenExpired
	}
	for _, perm := range permissions {
		found := false
		for _, p := range *result.Permissions {
			if *p.RSName == perm {
				found = true
				break
			}
		}
		if !found {
			return "", iam.ErrInsufficientPermissions
		}
	}
	part, err := base64.RawStdEncoding.DecodeString(strings.Split(token, ".")[1])
	if err != nil {
		return "", errors.Wrap(err, "base64")
	}
	var v struct {
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal([]byte(part), &v); err != nil {
		return "", err
	}
	if v.Sub == "" {
		panic(errors.New("malformed keycloak token passed retrospection"))
	}
	return v.Sub, nil
}
