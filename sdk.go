package main

import "strings"

type ScrimmageRewarderService interface {
}

type scrimmageRewarderServiceImpl struct {
	config *rewarderConfig
}

func InitScrimmageRewarderService(
	apiServerEndpoint string,
	privateKey string,
	namespace string,
	options ...RewarderOptionFnc,
) (ScrimmageRewarderService, error) {
	config := &rewarderConfig{
		apiServerEndpoint: apiServerEndpoint,
		privateKeys: map[string]string{
			"default": privateKey,
		},
		serviceMap: map[ServiceType]string{
			ServiceType_API: "api",
			ServiceType_P2E: "p2e",
			ServiceType_FED: "fed",
			ServiceType_NBC: "nbc",
		},
		namespace: namespace,
	}

	for _, runableOption := range options {
		runableOption(config)
	}

	if isUrlHasValidProtocol := validateURLProtocol(config.apiServerEndpoint, config.secure); !isUrlHasValidProtocol {
		return nil, ErrInvalidURLProtocol
	}

	config.apiServerEndpoint, _ = strings.CutSuffix(config.apiServerEndpoint, "/")

	return &scrimmageRewarderServiceImpl{
		config: config,
	}, nil
}
