package scrimmage

import "net/http"

type rewarderConfig struct {
	apiServerEndpoint         string
	privateKeys               map[string]string
	namespace                 string
	services                  []ServiceType
	logLevel                  LogLevel
	logger                    Logger
	secure                    bool
	validateAPIServerEndpoint bool
	httpClient                *http.Client
}
