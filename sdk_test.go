package scrimmage_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	scrimmage "github.com/Scrimmage-co/golang-sdk"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_SDK_InitOK(t *testing.T) {
	const (
		privateKey = "MOCK_PRIVATE_KEY"
		namespace  = "isolated-testing"
	)

	mockedScrimmageBackendHandler := gin.Default()
	mockedScrimmageBackendHandler.GET(":service/system/status", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"ok": true,
		})
	})

	mockedScrimmageBackendHandler.GET("api/rewarders/keys/@me", func(ctx *gin.Context) {
		assert.Equal(t, "Token "+privateKey, ctx.GetHeader("Authorization"))
		assert.Equal(t, namespace, ctx.GetHeader("Scrimmage-Namespace"))

		ctx.JSON(200, gin.H{
			"ok": true,
		})
	})

	mockedScrimmageBackendServer := httptest.NewServer(mockedScrimmageBackendHandler)
	apiServerEndpoint := mockedScrimmageBackendServer.URL

	_, err := scrimmage.InitRewarder(
		context.Background(),
		apiServerEndpoint,
		privateKey,
		namespace,
		scrimmage.WithSecure(false),
	)

	assert.NoError(t, err)
}

func Test_SDK_GetUserTokenForbidden(t *testing.T) {
	var (
		privateKey = "MOCK_PRIVATE_KEY"
		namespace  = "isolated-testing"

		userId         = "userId"
		tags           = []string{"tag-a", "tag-b"}
		userProperties = map[string]any{
			"name":    "user",
			"balance": 500000,
		}
	)

	mockedScrimmageBackendHandler := gin.Default()
	mockedScrimmageBackendHandler.POST("api/integrations/users", func(ctx *gin.Context) {
		ctx.JSON(403, gin.H{
			"message":    "Forbidden resource",
			"error":      "Forbidden",
			"statusCode": 403,
		})
	})

	mockedScrimmageBackendServer := httptest.NewServer(mockedScrimmageBackendHandler)
	apiServerEndpoint := mockedScrimmageBackendServer.URL

	sdk, err := scrimmage.InitRewarder(
		context.Background(),
		apiServerEndpoint,
		privateKey,
		namespace,
		scrimmage.WithSecure(false),
		scrimmage.WithValidateAPIServerEndpoint(false),
	)

	assert.NoError(t, err)

	_, err = sdk.User.GetUserToken(context.Background(), scrimmage.GetUserTokenRequest{
		UserID:     userId,
		Tags:       tags,
		Properties: userProperties,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, scrimmage.ErrForbidden)
}

func Test_SDK_GetUserTokenOK(t *testing.T) {
	var (
		privateKey = "MOCK_PRIVATE_KEY"
		namespace  = "isolated-testing"

		userId         = "userId"
		tags           = []string{"tag-a", "tag-b"}
		userProperties = map[string]any{
			"name":    "user",
			"balance": 500000,
		}

		userToken = "userToken"
	)

	mockedScrimmageBackendHandler := gin.Default()
	mockedScrimmageBackendHandler.POST("api/integrations/users", func(ctx *gin.Context) {
		var reqBody scrimmage.GetUserTokenRequest
		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.Abort()
			return
		}

		assert.Equal(t, userId, reqBody.UserID)
		assert.Subset(t, reqBody.Tags, tags)
		assert.EqualValues(t, userProperties["name"], reqBody.Properties["name"])
		assert.EqualValues(t, userProperties["balance"], reqBody.Properties["balance"])
		assert.Equal(t, "Token "+privateKey, ctx.GetHeader("Authorization"))
		assert.Equal(t, namespace, ctx.GetHeader("Scrimmage-Namespace"))

		ctx.JSON(200, scrimmage.GetUserTokenResponse{
			Token: userToken,
		})
	})

	mockedScrimmageBackendServer := httptest.NewServer(mockedScrimmageBackendHandler)
	apiServerEndpoint := mockedScrimmageBackendServer.URL

	sdk, err := scrimmage.InitRewarder(
		context.Background(),
		apiServerEndpoint,
		privateKey,
		namespace,
		scrimmage.WithSecure(false),
		scrimmage.WithValidateAPIServerEndpoint(false),
	)

	assert.NoError(t, err)

	userTokenResult, err := sdk.User.GetUserToken(context.Background(), scrimmage.GetUserTokenRequest{
		UserID:     userId,
		Tags:       tags,
		Properties: userProperties,
	})

	assert.NoError(t, err)
	assert.Equal(t, userToken, userTokenResult)
}

func Test_SDK_TrackRewardableOnceOK(t *testing.T) {
	var (
		privateKey = "MOCK_PRIVATE_KEY"
		namespace  = "isolated-testing"

		betEvent = scrimmage.BetEvent{
			BetType:     scrimmage.BetType_Single,
			IsLive:      false,
			Odds:        1.5,
			Description: "lorem ipsum",
			WagerAmount: 1000,
			NetProfit:   scrimmage.GetPtrOf[float64](500),
			Outcome:     scrimmage.GetPtrOf[scrimmage.BetOutcome]("win"),
			BetDate:     scrimmage.BetDate(time.Now().UnixMilli()),
			Bets: []scrimmage.SingleBet{
				{
					Type:           scrimmage.SingleBetType_Spread,
					Odds:           1.5,
					TeamBetOn:      scrimmage.GetPtrOf("team a"),
					TeamBetAgainst: scrimmage.GetPtrOf("team b"),
					League:         scrimmage.BetLeague("nba"),
					Sport:          scrimmage.BetSport("basketball"),
				},
			},
		}

		expectedRequestBodyInJson, _ = json.Marshal(scrimmage.CreateIntegrationRewardRequest{
			EventID:  scrimmage.GetPtrOf("uniqueEventId"),
			UserID:   "userId",
			DataType: scrimmage.BetDataType_BetExecuted,
			Body:     betEvent,
		})
	)

	mockedScrimmageBackendHandler := gin.Default()
	mockedScrimmageBackendHandler.GET(":service/system/status", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"ok": true,
		})
	})

	mockedScrimmageBackendHandler.GET("api/rewarders/keys/@me", func(ctx *gin.Context) {
		assert.Equal(t, "Token "+privateKey, ctx.GetHeader("Authorization"))
		assert.Equal(t, namespace, ctx.GetHeader("Scrimmage-Namespace"))

		ctx.JSON(200, gin.H{
			"ok": true,
		})
	})

	mockedScrimmageBackendHandler.POST("api/integrations/rewards", func(ctx *gin.Context) {
		rawRequestBody, err := io.ReadAll(ctx.Request.Body)
		assert.NoError(t, err)

		assert.JSONEq(t, string(rawRequestBody), string(expectedRequestBodyInJson))
		bodyJson, err := json.Marshal(betEvent)
		assert.NoError(t, err)

		ctx.JSON(200, scrimmage.CreateIntegrationRewardResponse{
			Namespace: namespace,
			EventID:   scrimmage.GetPtrOf("uniqueEventId"),
			DataType:  scrimmage.GetPtrOf(scrimmage.BetDataType_BetExecuted),
			Body:      bodyJson,
		})
	})

	mockedScrimmageBackendServer := httptest.NewServer(mockedScrimmageBackendHandler)
	apiServerEndpoint := mockedScrimmageBackendServer.URL

	sdk, err := scrimmage.InitRewarder(
		context.Background(),
		apiServerEndpoint,
		privateKey,
		namespace,
		scrimmage.WithSecure(false),
	)

	assert.NoError(t, err)

	_, err = sdk.Reward.TrackRewardableOnce(
		context.Background(),
		"userId",
		scrimmage.BetDataType_BetExecuted,
		scrimmage.GetPtrOf("uniqueEventId"), betEvent,
	)

	assert.NoError(t, err)
}

func Test_SDK_TrackRewardableMultipleDataOK(t *testing.T) {
	var (
		privateKey = "MOCK_PRIVATE_KEY"
		namespace  = "isolated-testing"

		betEvents = []scrimmage.BetEvent{
			{
				BetType:     scrimmage.BetType_Single,
				IsLive:      false,
				Odds:        1.5,
				Description: "lorem ipsum",
				WagerAmount: 1000,
				NetProfit:   scrimmage.GetPtrOf[float64](500),
				Outcome:     scrimmage.GetPtrOf[scrimmage.BetOutcome]("win"),
				BetDate:     scrimmage.BetDate(time.Now().UnixMilli()),
				Bets: []scrimmage.SingleBet{
					{
						Type:           scrimmage.SingleBetType_Spread,
						Odds:           1.5,
						TeamBetOn:      scrimmage.GetPtrOf("team a"),
						TeamBetAgainst: scrimmage.GetPtrOf("team b"),
						League:         scrimmage.BetLeague("nba"),
						Sport:          scrimmage.BetSport("basketball"),
					},
				},
			},
			{
				BetType:     scrimmage.BetType_Single,
				IsLive:      false,
				Odds:        1.5,
				Description: "lorem ipsum",
				WagerAmount: 1000,
				NetProfit:   scrimmage.GetPtrOf[float64](500),
				Outcome:     scrimmage.GetPtrOf[scrimmage.BetOutcome]("win"),
				BetDate:     scrimmage.BetDate(time.Now().UnixMilli()),
				Bets: []scrimmage.SingleBet{
					{
						Type:           scrimmage.SingleBetType_Spread,
						Odds:           1.5,
						TeamBetOn:      scrimmage.GetPtrOf("team a"),
						TeamBetAgainst: scrimmage.GetPtrOf("team b"),
						League:         scrimmage.BetLeague("nba"),
						Sport:          scrimmage.BetSport("basketball"),
					},
				},
			},
		}

		expectedRequestBodyInJson, _ = json.Marshal(scrimmage.CreateIntegrationRewardRequest{
			UserID:   "userId",
			DataType: scrimmage.BetDataType_BetExecuted,
			Body:     betEvents[0],
		})
	)

	mockedScrimmageBackendHandler := gin.Default()
	mockedScrimmageBackendHandler.GET(":service/system/status", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"ok": true,
		})
	})

	mockedScrimmageBackendHandler.GET("api/rewarders/keys/@me", func(ctx *gin.Context) {
		assert.Equal(t, "Token "+privateKey, ctx.GetHeader("Authorization"))
		assert.Equal(t, namespace, ctx.GetHeader("Scrimmage-Namespace"))

		ctx.JSON(200, gin.H{
			"ok": true,
		})
	})

	mockedScrimmageBackendHandler.POST("api/integrations/rewards", func(ctx *gin.Context) {
		rawRequestBody, err := io.ReadAll(ctx.Request.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, string(rawRequestBody), string(expectedRequestBodyInJson))

		bodyJson, err := json.Marshal(betEvents[0])
		assert.NoError(t, err)

		ctx.JSON(200, scrimmage.CreateIntegrationRewardResponse{
			Namespace: namespace,
			EventID:   scrimmage.GetPtrOf("uniqueEventId"),
			DataType:  scrimmage.GetPtrOf(scrimmage.BetDataType_BetExecuted),
			Body:      bodyJson,
		})
	})

	mockedScrimmageBackendServer := httptest.NewServer(mockedScrimmageBackendHandler)
	apiServerEndpoint := mockedScrimmageBackendServer.URL

	sdk, err := scrimmage.InitRewarder(
		context.Background(),
		apiServerEndpoint,
		privateKey,
		namespace,
		scrimmage.WithSecure(false),
	)

	assert.NoError(t, err)

	result, err := sdk.Reward.TrackRewardable(
		context.Background(),
		"userId",
		scrimmage.BetDataType_BetExecuted,
		betEvents[0],
		betEvents[1],
	)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
}
