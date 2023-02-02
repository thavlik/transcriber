package base

import (
	"context"
	"sync"
	"time"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type KeyCloakAdminOptions struct {
	BasicAuth
	Realm string
}

type KeyCloakOptions struct {
	Endpoint     string
	Admin        KeyCloakAdminOptions
	ClientID     string
	ClientSecret string
	Realm        string
}

func (o *KeyCloakOptions) IsSet() bool {
	return o.Endpoint != ""
}

func (o *KeyCloakOptions) Ensure() {
	if o.Endpoint == "" {
		panic(errors.New("missing --keycloak-endpoint"))
	}
	if o.Admin.Realm == "" {
		panic(errors.New("missing --keycloak-admin-realm"))
	}
	if o.Admin.Username == "" {
		panic(errors.New("missing --keycloak-admin-username"))
	}
	if o.Admin.Password == "" {
		panic(errors.New("missing --keycloak-admin-password"))
	}
	if o.ClientID == "" {
		panic(errors.New("missing --keycloak-client-id"))
	}
	if o.ClientSecret == "" {
		panic(errors.New("missing --keycloak-client-secret"))
	}
	if o.Realm == "" {
		panic(errors.New("missing --keycloak-realm"))
	}
}

func KeyCloakEnv(o *KeyCloakOptions, required bool) *KeyCloakOptions {
	if o == nil {
		o = &KeyCloakOptions{}
	}
	CheckEnv("KC_ENDPOINT", &o.Endpoint)
	CheckEnv("KC_ADMIN_USERNAME", &o.Admin.Username)
	CheckEnv("KC_ADMIN_PASSWORD", &o.Admin.Password)
	CheckEnv("KC_ADMIN_REALM", &o.Admin.Realm)
	CheckEnv("KC_CLIENT_ID", &o.ClientID)
	CheckEnv("KC_CLIENT_SECRET", &o.ClientSecret)
	CheckEnv("KC_REALM", &o.Realm)
	if required {
		o.Ensure()
	}
	return o
}

func AddKeyCloakFlags(cmd *cobra.Command, o *KeyCloakOptions) {
	cmd.PersistentFlags().StringVar(&o.Endpoint, "keycloak-endpoint", "", "keycloak http/https service endpoint")
	cmd.PersistentFlags().StringVar(&o.Admin.Username, "keycloak-admin-username", "", "keycloak admin username")
	cmd.PersistentFlags().StringVar(&o.Admin.Password, "keycloak-admin-password", "", "keycloak admin password")
	cmd.PersistentFlags().StringVar(&o.Admin.Realm, "keycloak-admin-realm", "master", "keycloak admin realm")
	cmd.PersistentFlags().StringVar(&o.ClientID, "keycloak-client-id", "", "keycloak client id")
	cmd.PersistentFlags().StringVar(&o.ClientSecret, "keycloak-client-secret", "", "keycloak client secret")
	cmd.PersistentFlags().StringVar(&o.Realm, "keycloak-realm", "", "keycloak application realm")
}

func ConnectKeyCloak(o *KeyCloakOptions) *KeyCloak {
	k := &KeyCloak{
		GoCloak:         gocloak.NewClient(o.Endpoint),
		KeyCloakOptions: *o,
	}
	if _, err := k.GetAccessToken(context.Background()); err != nil {
		panic(errors.Wrap(err, "failed to connect to keycloak"))
	}
	DefaultLog.Debug("connected to keycloak", Elapsed(start))
	return k
}

type KeyCloak struct {
	*gocloak.GoCloak
	KeyCloakOptions
	Token          *gocloak.JWT
	l              sync.Mutex
	tokenExpires   time.Time
	refreshExpires time.Time
}

func (k *KeyCloak) GetAccessToken(ctx context.Context) (string, error) {
	k.l.Lock()
	defer k.l.Unlock()
	t := time.Now()
	padSeconds := 10
	if t.Before(k.tokenExpires) {
		return k.Token.AccessToken, nil
	}
	var err error
	if k.Token, err = k.LoginAdmin(
		ctx,
		k.Admin.Username,
		k.Admin.Password,
		k.Admin.Realm,
	); err != nil {
		return "", errors.Wrap(err, "LoginAdmin")
	}
	k.tokenExpires = t.Add(time.Duration(k.Token.ExpiresIn-padSeconds) * time.Second)
	k.refreshExpires = t.Add(time.Duration(k.Token.RefreshExpiresIn-padSeconds) * time.Second)
	return k.Token.AccessToken, nil
}
