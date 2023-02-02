package keycloak

import (
	"github.com/Nerzal/gocloak/v12"
	"github.com/thavlik/transcriber/base/pkg/base"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

type keyCloakIAM struct {
	kc  *base.KeyCloak
	log *zap.Logger
}

func NewKeyCloakIAM(
	kc *base.KeyCloak,
	log *zap.Logger,
) iam.IAM {
	return &keyCloakIAM{kc, log}
}

func convertUser(user *gocloak.User) *iam.User {
	return &iam.User{
		ID:        gocloak.PString(user.ID),
		Username:  gocloak.PString(user.Username),
		Email:     gocloak.PString(user.Email),
		FirstName: gocloak.PString(user.FirstName),
		LastName:  gocloak.PString(user.LastName),
	}
}
func convertUsers(users []*gocloak.User) []*iam.User {
	var result []*iam.User
	for _, user := range users {
		result = append(result, convertUser(user))
	}
	return result
}
