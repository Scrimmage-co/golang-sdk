package scrimmage_test

import (
	"context"
	"net/http/httptest"
	"testing"

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
		isAuthHeaderValid := ctx.GetHeader("Authorization") == "Token "+privateKey
		isNamespaceValid := ctx.GetHeader("Scrimmage-Namespace") == namespace

		if !isAuthHeaderValid || !isNamespaceValid {
			ctx.Abort()
			return
		}

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
