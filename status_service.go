package scrimmage

import (
	"context"
)

type statusCheckResult struct {
	serviceType ServiceType
	ok          bool
	err         error
}

type StatusService interface {
	Verify(ctx context.Context)
}

type statusServiceImpl struct {
	api    API
	config *rewarderConfig
}

func newStatusService(
	api API,
	config *rewarderConfig,
) StatusService {
	return &statusServiceImpl{
		api:    api,
		config: config,
	}
}

func (s *statusServiceImpl) Verify(ctx context.Context) {
	if !s.config.validateAPIServerEndpoint {
		return
	}

	if !s.getOverallServiceStatus(ctx) {
		s.config.logger.Error("Rewarder API is not available")
	}

	if err := s.api.GetRewarderKeyDetails(ctx); err != nil {
		s.config.logger.Error("Rewarder API key is invalid")
	}

}

func (s *statusServiceImpl) getOverallServiceStatus(ctx context.Context) bool {
	ch := make(chan statusCheckResult, len(s.config.services))
	for _, service := range s.config.services {
		go func(ctx context.Context, ch chan<- statusCheckResult, serviceType ServiceType) {
			result := statusCheckResult{
				serviceType: serviceType,
			}

			result.ok = true
			if err := s.api.GetServiceStatus(ctx, serviceType); err != nil {
				result.ok = false
				result.err = err
			}

			ch <- result
		}(ctx, ch, service)
	}

	var isAnyServiceDown bool
	for i := 0; i < len(s.config.services); i++ {
		serviceStatus := <-ch
		if !serviceStatus.ok {
			isAnyServiceDown = true
		}
	}

	return !isAnyServiceDown
}
