package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xuxant/Darwin_API/pkg/aws"
	"github.com/xuxant/Darwin_API/pkg/dto"
	"github.com/xuxant/Darwin_API/pkg/logs"
	"github.com/xuxant/Darwin_API/pkg/models"
	"github.com/xuxant/Darwin_API/pkg/utils"
	"gorm.io/gorm"
	"net/http"
	"strings"
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

	account, err := aws.ValidateCredentials(registerCloud.AccessKey, registerCloud.SecretAccessKey)
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	registerCloud.UserId = c.GetString("user")

	if err := h.svc.AddCloudCredentials(registerCloud); err != nil {
		restErr := utils.NewInternalServerError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}

	response := models.ResponseBody{
		Success: true,
		Data: models.CloudRegister{
			Provider: registerCloud.Provider,
			Identity: *account,
		},
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetAllObjectsFromBucket(c *gin.Context) {
	bucketName := c.Param("bucketName")
	creds, err := h.svc.FetchCloudCredentials(c.GetString("user"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			restErr := utils.NewNotFoundError(err.Error())
			c.JSON(restErr.Status, restErr)
			return
		}
		restErr := utils.NewInternalServerError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}

	objects, err := aws.GetAllObjects(creds.AccessKey, creds.SecretAccessKey, bucketName)
	if err != nil {
		restErr := utils.NewUnauthorizedError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	c.JSON(http.StatusOK, models.ResponseBody{
		Success: true,
		Data:    objects,
	})

}

func (h *Handler) GetAllBuckets(c *gin.Context) {
	userId := c.GetString("user")
	creds, err := h.svc.FetchCloudCredentials(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			restErr := utils.NewNotFoundError(err.Error())
			c.JSON(restErr.Status, restErr)
			return
		}
		restErr := utils.NewInternalServerError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	buckets, err := aws.GetAllBuckets(creds.AccessKey, creds.SecretAccessKey)
	if err != nil {
		restErr := utils.NewUnauthorizedError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}

	response := models.ResponseBody{
		Success: true,
		Data:    buckets,
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) MountDataObject(c *gin.Context) {
	var dataMetadata []models.DataMetadata
	var dataObject models.AddS3Object

	if err := c.ShouldBindJSON(&dataObject); err != nil {
		restErr := utils.NewBadRequestError("Invalid Json Body")
		c.JSON(restErr.Status, restErr)
		return
	}

	uId := c.GetString("user")
	for _, objects := range dataObject.Objects {
		fileName := strings.Split(objects, "/")
		dataMetadata = append(dataMetadata, models.DataMetadata{
			DataId:     utils.GenerateUUID(),
			UserId:     uId,
			BucketName: dataObject.BucketName,
			FileName:   fileName[len(fileName)-1],
			ObjectKey:  objects,
		})
	}
	if err := h.svc.AddS3Objects(dataMetadata); err != nil {
		restErr := utils.NewInternalServerError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	c.JSON(http.StatusCreated, models.ResponseBody{
		Success: true,
	})
}

func (h *Handler) FetchMountedFiles(c *gin.Context) {
	results, err := h.svc.FetchMountedFiles(c.GetString("user"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			restErr := utils.NewNotFoundError(err.Error())
			c.JSON(restErr.Status, restErr)
			return
		}
		restErr := utils.NewInternalServerError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	var mountedFiles []models.MountedFiles
	for _, result := range results {
		mountedFiles = append(mountedFiles, models.MountedFiles{
			DataId:   result.DataId,
			FileName: result.FileName,
		})
	}
	c.JSON(http.StatusOK, models.ResponseBody{
		Success: true,
		Data:    mountedFiles,
	})
}

func (h *Handler) UpdateCloudCredentials(c *gin.Context) {
	var registerCloud models.CloudCredentials
	logger := logs.LoggerForRequest(c)
	logger.Info().Msg("Registering AWS")

	if err := c.ShouldBindJSON(&registerCloud); err != nil {
		restErr := utils.NewBadRequestError("Invalid Json Body")
		c.JSON(restErr.Status, restErr)
		return
	}

	_, err := aws.ValidateCredentials(registerCloud.AccessKey, registerCloud.SecretAccessKey)
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	err = h.svc.UpdateCredentials(c.GetString("user"), registerCloud.AccessKey, registerCloud.SecretAccessKey)
	if err != nil {
		restErr := utils.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	c.JSON(http.StatusOK, models.ResponseBody{
		Success: true,
	})
}
