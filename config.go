package main

import "net/http"

type rewarderConfig struct {
	apiServerEndpoint         string
	privateKeys               map[string]string
	namespace                 string
	serviceMap                map[ServiceType]string
	logLevel                  LogLevel
	logger                    Logger
	secure                    bool
	validateAPIServerEndpoint bool
	httpClient                *http.Client
}
