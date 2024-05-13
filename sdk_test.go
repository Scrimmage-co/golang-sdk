package scrimmage_test

import (
	"context"
	"net/http/httptest"
	"testing"

	scrimmage "github.com/Scrimmage-co/golang-sdk"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_InitSDK_OK(t *testing.T) {
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
