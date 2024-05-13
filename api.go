package scrimmage

import (
	"context"
	"fmt"
	"net/http"
)

const (
	SERVICE_STATUS_PATH      = "system/status"
	REWARDER_KEY_DETAIL_PATH = "rewarders/keys/@me"
)

type API interface {
	GetServiceStatus(ctx context.Context, serviceName ServiceType) error
	GetRewarderKeyDetails(ctx context.Context) error
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

	r, err := a.config.httpClient.Get(finalUrl)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return ErrStatusCodeIsNotOK
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

	if _, err := a.config.httpClient.Do(req); err != nil {
		return err
	}

	return nil
}
