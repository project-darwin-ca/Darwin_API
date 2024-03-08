package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xuxant/Darwin_API/pkg/aws"
	"github.com/xuxant/Darwin_API/pkg/logs"
	"github.com/xuxant/Darwin_API/pkg/models"
	"github.com/xuxant/Darwin_API/pkg/utils"
	"net/http"
)

func RegisterAWS(c *gin.Context) {
	var registerCloud models.RegisterCloudProvider
	logger := logs.LoggerForRequest(c)
	logger.Info().Msg("Registering AWS")

	if err := c.ShouldBindJSON(&registerCloud); err != nil {
		restErr := utils.NewBadRequestError("Invalid Json Body")
		c.JSON(restErr.Status, restErr)
		return
	}

	if err := aws.ValidateCredentials(registerCloud.AccessKey, registerCloud.SecretAccessKey); err != nil {
		restErr := utils.NewBadRequestError("invalid access key and secret key")
		c.JSON(restErr.Status, restErr)
		return
	}

	c.JSON(http.StatusCreated, registerCloud)
}
