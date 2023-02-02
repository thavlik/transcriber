package keycloak

import (
	"context"

	"github.com/pkg/errors"
)

func (i *keyCloakIAM) refreshAccessToken(ctx context.Context) (string, error) {
	v, err := i.kc.GetAccessToken(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to refresh admin access token")
	}
	return v, nil
}
