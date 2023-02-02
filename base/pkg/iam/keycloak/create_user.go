package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

func (i *keyCloakIAM) CreateUser(
	user *iam.User,
	password string,
	temporary bool,
) (string, error) {
	accessToken, err := i.refreshAccessToken(context.Background())
	if err != nil {
		return "", err
	}
	i.log.Debug("creating user",
		zap.String("realm", i.kc.Realm),
		zap.String("user.ID", user.ID),
		zap.String("user.Username", user.Username),
		zap.String("user.FirstName", user.FirstName),
		zap.String("user.LastName", user.LastName),
		zap.String("user.Email", user.Email))
	user.ID, err = i.kc.CreateUser(
		context.Background(),
		accessToken,
		i.kc.Realm,
		gocloak.User{
			ID:            gocloak.StringP(user.ID),
			Username:      gocloak.StringP(user.Username),
			FirstName:     gocloak.StringP(user.FirstName),
			LastName:      gocloak.StringP(user.LastName),
			Email:         gocloak.StringP(user.Email),
			Enabled:       gocloak.BoolP(user.Enabled),
			EmailVerified: gocloak.BoolP(true),
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "keycloak create user")
	}
	if err := i.kc.SetPassword(
		context.Background(),
		accessToken,
		user.ID,
		i.kc.Realm,
		password,
		temporary,
	); err != nil {
		return "", errors.Wrap(err, "keycloak set password")
	}
	i.log.Debug("created user",
		zap.String("realm", i.kc.Realm),
		zap.String("user.ID", user.ID),
		zap.String("user.Username", user.Username),
		zap.String("user.FirstName", user.FirstName),
		zap.String("user.LastName", user.LastName),
		zap.String("user.Email", user.Email))
	return user.ID, nil
}
