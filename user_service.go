package scrimmage

import (
	"context"
)

type GetUserTokenRequest struct {
	UserID     string         `json:"id"`
	Tags       []string       `json:"tags"`
	Properties map[string]any `json:"properties"`
}

type UserService interface {
	GetUserToken(ctx context.Context, payload GetUserTokenRequest) (string, error)
}

type userServiceImpl struct {
	config *rewarderConfig
	api    API
}

func newUserServiceImpl(
	config *rewarderConfig,
	api API,
) UserService {
	return &userServiceImpl{
		config: config,
		api:    api,
	}
}

func (u *userServiceImpl) GetUserToken(ctx context.Context, payload GetUserTokenRequest) (string, error) {
	return u.api.GetUserToken(ctx, payload)
}
