package scrimmage

import "net/http"

type ServiceType string

const (
	ServiceType_API ServiceType = "api"
	ServiceType_P2E ServiceType = "p2e"
	ServiceType_FED ServiceType = "fed"
	ServiceType_NBC ServiceType = "nbc"
)

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
