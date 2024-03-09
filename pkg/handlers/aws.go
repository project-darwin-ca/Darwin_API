package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuxant/Darwin_API/pkg/aws"
	"github.com/xuxant/Darwin_API/pkg/dto"
	"github.com/xuxant/Darwin_API/pkg/logs"
	"github.com/xuxant/Darwin_API/pkg/models"
	"github.com/xuxant/Darwin_API/pkg/utils"
	"gorm.io/gorm"
	"net/http"
)

type Handler struct {
	svc *dto.Service
}

func NewDataHandler(db *gorm.DB) *Handler {
	svc := dto.NewDtoService(db)
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterAWS(c *gin.Context) {
	var registerCloud models.CloudCredentials
	logger := logs.LoggerForRequest(c)
	logger.Info().Msg("Registering AWS")

	if err := c.ShouldBindJSON(&registerCloud); err != nil {
		restErr := utils.NewBadRequestError("Invalid Json Body")
		c.JSON(restErr.Status, restErr)
		return
	}
	fmt.Println(registerCloud)

	account, err := aws.ValidateCredentials(registerCloud.AccessKey, registerCloud.SecretAccessKey)
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}

	if err := h.svc.AddCloudCredentials(registerCloud); err != nil {
		restErr := utils.NewInternalServerError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}

	var response models.ResponseBody
	response.Success = true
	response.Data.Identity = *account
	response.Data.Provider = registerCloud.Provider

	c.JSON(http.StatusCreated, response)
}
