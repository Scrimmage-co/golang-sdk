package scrimmage

type RewarderOptionFnc func(config *rewarderConfig) *rewarderConfig

func WithLogLevel(logLevel LogLevel) RewarderOptionFnc {
	return func(config *rewarderConfig) *rewarderConfig {
		config.logLevel = logLevel
		return config
	}
}

func WithLogger(logger Logger) RewarderOptionFnc {
	return func(config *rewarderConfig) *rewarderConfig {
		config.logger = logger
		return config
	}
}

func WithSecure(secure bool) RewarderOptionFnc {
	return func(config *rewarderConfig) *rewarderConfig {
		config.secure = secure
		return config
	}
}

func WithValidateAPIServerEndpoint(isValidate bool) RewarderOptionFnc {
	return func(config *rewarderConfig) *rewarderConfig {
		config.validateAPIServerEndpoint = isValidate
		return config
	}
}
