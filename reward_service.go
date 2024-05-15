package scrimmage

import (
	"context"
	"sync"
)

type RewardService interface {
	TrackRewardable(ctx context.Context, userId string, dataType BetDataType, event ...BetEvent) ([]CreateIntegrationRewardResponse, error)
	TrackRewardableOnce(ctx context.Context, userId string, dataType BetDataType, eventID *string, event BetEvent) (CreateIntegrationRewardResponse, error)
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

func (s *rewardServiceImpl) TrackRewardable(ctx context.Context, userId string, dataType BetDataType, events ...BetEvent) ([]CreateIntegrationRewardResponse, error) {
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
			event BetEvent,
		) {
			defer wg.Done()

			result, _ := s.api.CreateIntegrationReward(ctx, CreateIntegrationRewardRequest{
				UserID:   userId,
				DataType: dataType,
				Body:     event,
			})

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

func (s *rewardServiceImpl) TrackRewardableOnce(ctx context.Context, userId string, dataType BetDataType, eventID *string, event BetEvent) (CreateIntegrationRewardResponse, error) {
	return s.api.CreateIntegrationReward(ctx, CreateIntegrationRewardRequest{
		UserID:   userId,
		EventID:  eventID,
		DataType: dataType,
		Body:     event,
	})
}
