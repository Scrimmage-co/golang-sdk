package scrimmage

import (
	"context"
	"sync"
)

type RewardService interface {
	TrackRewardable(ctx context.Context, userId string, dataType string, event ...interface{}) ([]CreateIntegrationRewardResponse, error)
	TrackRewardableOnce(ctx context.Context, userId string, dataType string, eventID *string, events interface{}) (CreateIntegrationRewardResponse, error)
}

type rewardServiceImpl struct {
	config *rewarderConfig
	api    API
}

func newRewardService(
	config *rewarderConfig,
	api API,
) RewardService {
	return &rewardServiceImpl{
		config: config,
		api:    api,
	}
}

func (s *rewardServiceImpl) TrackRewardable(ctx context.Context, userId string, dataType string, events ...interface{}) ([]CreateIntegrationRewardResponse, error) {
	var (
		mutex sync.Mutex
		wg    sync.WaitGroup
	)

	wg.Add(len(events))
	results := make([]CreateIntegrationRewardResponse, 0, len(events))

	for _, event := range events {
		go func(
			mutex *sync.Mutex,
			wg *sync.WaitGroup,
			event interface{},
		) {
			defer wg.Done()

			result, err := s.api.CreateIntegrationReward(ctx, CreateIntegrationRewardRequest{
				UserID:   userId,
				DataType: dataType,
				Body:     event,
			})

			if err != nil {
				s.config.logger.Warn("failed to call the integration reward : ", err)
			}

			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()
		}(
			&mutex,
			&wg,
			event,
		)
	}

	wg.Wait()

	return results, nil
}

func (s *rewardServiceImpl) TrackRewardableOnce(ctx context.Context, userId string, dataType string, eventID *string, events interface{}) (CreateIntegrationRewardResponse, error) {
	return s.api.CreateIntegrationReward(ctx, CreateIntegrationRewardRequest{
		UserID:   userId,
		EventID:  eventID,
		DataType: dataType,
		Body:     events,
	})
}
