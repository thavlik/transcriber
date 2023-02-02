package cognito

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/transcriber/base/pkg/iam"
	"go.uber.org/zap"
)

func (i *cognitoIAM) Login(
	ctx context.Context,
	username string,
	password string,
) (string, error) {
	var r struct {
		AuthParameters struct {
			Username string `json:"USERNAME"`
			Password string `json:"PASSWORD"`
		}
		AuthFlow string
		ClientId string
	}
	r.AuthParameters.Username = username
	r.AuthParameters.Password = password
	r.AuthFlow = "USER_PASSWORD_AUTH"
	r.ClientId = i.clientID
	body, err := json.Marshal(&r)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/", i.region)
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-amz-json-1.1")
	req.Header.Set("X-Amz-Target", "AWSCognitoIdentityProviderService.InitiateAuth")
	req = req.WithContext(ctx)
	resp, err := (&http.Client{
		Timeout: 15 * time.Second,
	}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		i.log.Error("failed cognito login", zap.String("err", string(body)))
		if resp.StatusCode == http.StatusBadRequest {
			return "", iam.ErrInvalidCredentials
		}
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}
	var result struct {
		ChallengeName        string
		AuthenticationResult struct {
			AccessToken  string
			IdToken      string
			RefreshToken string
		}
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.ChallengeName != "" {
		return "", errors.New("challenges not supported")
	}
	token := buildAuthBearerToken(
		result.AuthenticationResult.AccessToken,
		result.AuthenticationResult.IdToken,
		result.AuthenticationResult.RefreshToken,
		i.region,
		i.userPoolID,
		i.clientID,
	)
	return "Bearer " + token, nil
}

func buildAuthBearerToken(
	accessToken,
	idToken,
	refreshToken,
	region,
	userPoolId,
	clientId string,
) string {
	jwt, err := json.Marshal(&authHeader{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		IDToken:       idToken,
		CognitoRegion: region,
		UserPoolID:    userPoolId,
		ClientID:      clientId,
	})
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(jwt)
}
