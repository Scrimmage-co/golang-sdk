package scrimmage

type loggerService struct {
	config *rewarderConfig
}

func newLoggerService(
	config *rewarderConfig,
) *loggerService {
	return &loggerService{
		config: config,
	}
}

func (l *loggerService) Log(args ...interface{}) {
	if LogLevel_Log > l.config.logLevel {
		return
	}

	l.config.logger.Info(args...)
}

func (l *loggerService) Warn(args ...interface{}) {
	if LogLevel_Warn > l.config.logLevel {
		return
	}

	l.config.logger.Warn(args...)
}

func (l *loggerService) Debug(args ...interface{}) {
	if LogLevel_Debug > l.config.logLevel {
		return
	}

	l.config.logger.Debug(args...)
}

func (l *loggerService) Info(args ...interface{}) {
	if LogLevel_Info > l.config.logLevel {
		return
	}

	l.config.logger.Info(args...)
}

func (l *loggerService) Error(args ...interface{}) {
	if LogLevel_Error > l.config.logLevel {
		return
	}

	l.config.logger.Error(args...)
}
