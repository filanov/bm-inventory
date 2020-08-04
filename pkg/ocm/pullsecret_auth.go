package ocm

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type OCMAuthentication interface {
	AuthenticatePullSecret(ctx context.Context, pullSecret string) (userName string, err error)
}

type authentication service

var _ OCMAuthentication = &authentication{}

func (a authentication) AuthenticatePullSecret(ctx context.Context, pullSecret string) (userName string, err error) {

	//TODO cache pullSecret <-> Username
	con := a.client.connection
	request := con.Post()
	request.Path("/api/accounts_mgmt/v1/token_authorization")

	type TokenAuthorizationRequest struct {
		AuthorizationToken string `json:"authorization_token"`
	}

	tokenAuthorizationRequest := TokenAuthorizationRequest{
		AuthorizationToken: pullSecret,
	}

	var jsonData []byte
	jsonData, err = json.Marshal(tokenAuthorizationRequest)
	if err != nil {
		return "", err
	}
	request.Bytes(jsonData)

	postResp, err := request.SendContext(ctx)
	if err != nil || postResp.Status() != 200 {
		return "", err
	}

	type TokenAuthorizationResponse struct {
		Items []struct {
			Username string `json:"username"`
		} `json:"items"`
	}
	
	var tokenAuthorizationResponse TokenAuthorizationResponse
	if err := json.Unmarshal(postResp.Bytes(), &tokenAuthorizationResponse); err != nil {
		return "", err
	}
	logrus.Error(tokenAuthorizationResponse.Items[0].Username)
	return tokenAuthorizationResponse.Items[0].Username, nil
}
