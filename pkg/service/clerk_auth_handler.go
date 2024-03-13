package service

import (
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gin-gonic/gin"
	"github.com/xuxant/Darwin_API/pkg/utils"
	"strings"
)

func authHandlerMiddleware(client *jwks.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		sessionClaim, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{
			Token:      sessionToken,
			JWKSClient: client,
		})
		if err != nil {
			restError := utils.NewUnauthorizedError(err.Error())
			c.JSON(restError.Status, restError)
			c.Abort()
			return
		}
		c.Set("user", sessionClaim.Subject)

		c.Next()
	}
}
