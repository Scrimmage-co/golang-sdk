package scrimmage

import (
	"context"
	"net/http"
	"time"
)

type ScrimmageRewarder struct {
	config *rewarderConfig
	logger *loggerService
	status StatusService

	User   UserService
	Reward RewardService
}

func InitRewarder(
	ctx context.Context,
	apiServerEndpoint string,
	privateKey string,
	namespace string,
	options ...RewarderOptionFnc,
) (*ScrimmageRewarder, error) {
	sdk := &ScrimmageRewarder{}

	if err := sdk.setConfig(apiServerEndpoint, privateKey, namespace, options...); err != nil {
		return nil, err
	}

	apiClient := newAPI(sdk.config)

	sdk.status = newStatusService(apiClient, sdk.config)
	sdk.status.Verify(ctx)

	sdk.User = newUserServiceImpl(sdk.config, apiClient)
	sdk.Reward = newRewardService(sdk.config, apiClient)

	sdk.logger = newLoggerService(sdk.config)
	sdk.logger.Info("Rewarder Initiated")

	return sdk, nil
}

func (s *ScrimmageRewarder) setConfig(
	apiServerEndpoint string,
	privateKey string,
	namespace string,
	options ...RewarderOptionFnc,
) error {
	config := &rewarderConfig{
		apiServerEndpoint: apiServerEndpoint,
		privateKeys: map[string]string{
			"default": privateKey,
		},
		services:  []ServiceType{ServiceType_API, ServiceType_P2E, ServiceType_FED, ServiceType_NBC},
		namespace: namespace,
		httpClient: &RetryClient{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
			maxRetries: 3,
			baseDelay:  100 * time.Millisecond,
			maxBackoff: 1 * time.Second,
		},
		logLevel:                  LogLevel_Debug,
		logger:                    newDefaultLogger(),
		secure:                    true,
		validateAPIServerEndpoint: true,
	}

	for _, runableOption := range options {
		runableOption(config)
	}

	if isUrlHasValidProtocol := validateURLProtocol(config.apiServerEndpoint, config.secure); !isUrlHasValidProtocol {
		return ErrInvalidURLProtocol
	}

	config.apiServerEndpoint = CutSuffix(config.apiServerEndpoint, "/")
	s.config = config

	return nil
}
