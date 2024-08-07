package scrimmage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	SERVICE_STATUS_PATH            = "system/status"
	REWARDER_KEY_DETAIL_PATH       = "rewarders/keys/@me"
	GET_USER_TOKEN_PATH            = "integrations/users"
	CREATE_INTEGRATION_REWARD_PATH = "integrations/rewards"
)

type API interface {
	GetServiceStatus(ctx context.Context, serviceName ServiceType) error
	GetRewarderKeyDetails(ctx context.Context) error
	GetUserToken(ctx context.Context, payload GetUserTokenRequest) (string, error)
	CreateIntegrationReward(ctx context.Context, payload CreateIntegrationRewardRequest) (CreateIntegrationRewardResponse, error)
}

type apiImpl struct {
	config *rewarderConfig
}

func newAPI(config *rewarderConfig) API {
	return &apiImpl{
		config: config,
	}
}

func (a *apiImpl) GetServiceStatus(ctx context.Context, serviceName ServiceType) error {
	finalUrl := fmt.Sprintf("%s/%s/%s", a.config.apiServerEndpoint, serviceName, SERVICE_STATUS_PATH)

	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return err
	}

	res, err := a.config.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return ErrForbidden
	}

	if res.StatusCode == http.StatusBadRequest {
		var errorResponse BadRequestError
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			return err
		}

		return &errorResponse
	}

	return nil
}

func (a *apiImpl) GetRewarderKeyDetails(ctx context.Context) error {
	finalUrl := fmt.Sprintf("%s/%s/%s", a.config.apiServerEndpoint, ServiceType_API, REWARDER_KEY_DETAIL_PATH)

	req, err := http.NewRequest("GET", finalUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Token "+a.config.privateKeys["default"])
	req.Header.Set("Scrimmage-Namespace", a.config.namespace)
	req.Header.Set("Content-Type", "application/json")

	res, err := a.config.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return ErrForbidden
	}

	if res.StatusCode == http.StatusBadRequest {
		var errorResponse BadRequestError
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			return err
		}

		return &errorResponse
	}

	return nil
}

func (a *apiImpl) GetUserToken(ctx context.Context, payload GetUserTokenRequest) (string, error) {
	finalUrl := fmt.Sprintf("%s/%s/%s", a.config.apiServerEndpoint, ServiceType_API, GET_USER_TOKEN_PATH)

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", finalUrl, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+a.config.privateKeys["default"])
	req.Header.Set("Scrimmage-Namespace", a.config.namespace)
	req.Header.Set("Content-Type", "application/json")

	res, err := a.config.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return "", ErrForbidden
	}

	if res.StatusCode == http.StatusBadRequest {
		var errorResponse BadRequestError
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			return "", err
		}

		return "", &errorResponse
	}

	var responseBody GetUserTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return "", err
	}

	return responseBody.Token, nil
}

func (a *apiImpl) CreateIntegrationReward(ctx context.Context, payload CreateIntegrationRewardRequest) (CreateIntegrationRewardResponse, error) {
	finalUrl := fmt.Sprintf("%s/%s/%s", a.config.apiServerEndpoint, ServiceType_API, CREATE_INTEGRATION_REWARD_PATH)

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return CreateIntegrationRewardResponse{}, err
	}

	req, err := http.NewRequest("POST", finalUrl, bytes.NewReader(reqBody))
	if err != nil {
		return CreateIntegrationRewardResponse{}, err
	}

	req.Header.Set("Authorization", "Token "+a.config.privateKeys["default"])
	req.Header.Set("Scrimmage-Namespace", a.config.namespace)
	req.Header.Set("Content-Type", "application/json")

	res, err := a.config.httpClient.Do(req)
	if err != nil {
		return CreateIntegrationRewardResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return CreateIntegrationRewardResponse{}, ErrAccountIsNotLinked
	}

	if res.StatusCode == http.StatusForbidden {
		return CreateIntegrationRewardResponse{}, ErrForbidden
	}

	if res.StatusCode == http.StatusBadRequest {
		var errorResponse BadRequestError
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			return CreateIntegrationRewardResponse{}, err
		}

		return CreateIntegrationRewardResponse{}, &errorResponse
	}

	var responseBody CreateIntegrationRewardResponse
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		return CreateIntegrationRewardResponse{}, err
	}

	return responseBody, nil
}
