package cognito

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
)

// accessTokenClaims describes the fields we need to look at to validate an access token
type accessTokenClaims struct {
	Username      string   `json:"username"`
	TokenUse      string   `json:"token_use"`
	ClientID      string   `json:"client_id"`
	CognitoGroups []string `json:"cognito:groups"`
	jwt.StandardClaims
}

// idTokenClaims describes the fields we need to look at to validate an ID token
type idTokenClaims struct {
	TokenUse string `json:"token_use"`
	Realm    string `json:"custom:realm"`
	jwt.StandardClaims
}

type authHeader struct {
	AccessToken   string `json:"accessToken"`
	RefreshToken  string `json:"refreshToken"`
	IDToken       string `json:"idToken"`
	CognitoRegion string `json:"cognitoRegion"`
	UserPoolID    string `json:"userPoolId"`
	ClientID      string `json:"clientId"`
}

func retrieveAuthHeader(r string) (*authHeader, error) {
	split := strings.Split(r, "Bearer ")
	if len(split) != 2 {
		return nil, errors.New("basic authorization header missing Bearer prefix")
	}
	authJSON, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return nil, errors.Wrap(err, "failed to base64 decode auth header")
	}
	var header authHeader
	if err := json.Unmarshal(
		[]byte(authJSON),
		&header,
	); err != nil {
		return nil, errors.Wrap(err, "failed to decode auth header json")
	}
	return &header, nil
}

func (i *cognitoIAM) Authorize(
	ctx context.Context,
	token string,
	permissions []string,
) (string, error) {
	header, err := retrieveAuthHeader(token)
	if err != nil {
		return "", errors.Wrap(err, "retrieveAuthHeader")
	}

	// get the pubkey for the user pool
	region := aws.StringValue(i.cognito.Config.Region)
	url := "https://cognito-idp." + region + ".amazonaws.com/" + i.userPoolID + "/.well-known/jwks.json"
	awsPubKey, err := jwk.Fetch(ctx, url)
	if err != nil {
		return "", errors.Wrap(err, "cognito")
	}

	// parse access token
	accessToken, err := jwt.ParseWithClaims(
		header.AccessToken,
		&accessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// verify token is signed with RS256
			if token.Method != jwt.SigningMethodRS256 {
				return nil, fmt.Errorf("unexpected signing method for access token: %v", token.Header["alg"])
			}
			// make sure the token key id is contained in the relevant Cognito pool's pubkey list
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("access token has no header")
			}
			key, ok := awsPubKey.LookupKeyID(kid)
			if !ok {
				return nil, errors.New("access token key id not found in aws pubkey list")
			}
			// return the AWS key to Parse to verify the signature
			var raw interface{}
			return raw, key.Raw(&raw)
		})
	if err != nil {
		e, ok := err.(*jwt.ValidationError)
		if !i.allowTokenUseBeforeIssue || (ok && e.Errors&jwt.ValidationErrorIssuedAt == 0) { // Don't report error that token used before issued.
			return "", errors.Wrap(err, "parse access token")
		}
		accessToken.Valid = true
	}

	if !accessToken.Valid {
		return "", iam.ErrTokenExpired
	}

	// validate access token
	access, ok := accessToken.Claims.(*accessTokenClaims)
	if ok {
		// make sure a username is present
		if len(access.Username) == 0 {
			return "", errors.New("access token does not contain a username")
		}

		// verify token client ID is correct
		if access.ClientID != header.ClientID {
			return "", errors.New("access token has bad client id")
		}

		// verify token was issued from the correct user pool
		if access.Issuer != "https://cognito-idp."+i.region+".amazonaws.com/"+i.userPoolID {
			return "", errors.New("access token has bad issuer")
		}

		// verify token use is set to access
		if access.TokenUse != "access" {
			return "", errors.New("access token has bad use")
		}
	} else {
		return "", errors.New("access token claims are not valid")
	}

	// parse ID token
	idToken, err := jwt.ParseWithClaims(
		header.IDToken,
		&idTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// verify token is signed with RS256
			if token.Method != jwt.SigningMethodRS256 {
				return nil, fmt.Errorf("unexpected signing method for id token: %v", token.Header["alg"])
			}

			// make sure the token key id is contained in the relevant Cognito pool's pubkey list
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("id token has no header")
			}

			key, ok := awsPubKey.LookupKeyID(kid)
			if !ok {
				return nil, errors.New("id token key id not found in AWS pubkey list")
			}

			// return the AWS key to Parse to verify the signature
			var raw interface{}
			return raw, key.Raw(&raw)
		})
	if err != nil {
		//https://github.com/dgrijalva/jwt-go/issues/314#issuecomment-651329500
		e, ok := err.(*jwt.ValidationError)
		if !i.allowTokenUseBeforeIssue || (ok && e.Errors&jwt.ValidationErrorIssuedAt == 0) { // Don't report error that token used before issued?
			return "", errors.Wrap(err, "failed to parse id token")
		}
		idToken.Valid = true
	}

	// check time claims (iat, nbf, exp)
	if !idToken.Valid {
		return "", errors.New("id token is expired")
	}

	// validate ID token
	id, ok := idToken.Claims.(*idTokenClaims)
	if !ok {
		return "", errors.New("id token claims are not valid")
	}

	// verify token client ID is correct
	if id.Audience != header.ClientID {
		return "", errors.New("id token has bad client aud")
	}

	// verify token was issued from the correct user pool
	if id.Issuer != "https://cognito-idp."+i.region+".amazonaws.com/"+i.userPoolID {
		return "", errors.New("id token has bad issuer")
	}

	// verify token use is set to access
	if id.TokenUse != "id" {
		return "", errors.New("id token has bad use")
	}

	// Tokens are good. Ensure permissions are good.
	for _, perm := range permissions {
		found := false
		for _, p := range access.CognitoGroups {
			if p == perm {
				found = true
				break
			}
		}
		if !found {
			return "", iam.ErrInsufficientPermissions
		}
	}

	return access.Username, nil
}
