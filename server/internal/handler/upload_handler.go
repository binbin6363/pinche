package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pinche/internal/middleware"
	"pinche/internal/model"
	"pinche/internal/service"
)

type UploadHandler struct {
	service     *service.UploadService
	userService *service.UserService
}

func NewUploadHandler(uploadService *service.UploadService, userService *service.UserService) *UploadHandler {
	return &UploadHandler{
		service:     uploadService,
		userService: userService,
	}
}

// UploadImage handles POST /api/upload/image
// Query param: biz_type - "images" for chat images, "avatar" for user avatar
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "请选择文件"))
		return
	}

	bizType := c.Query("biz_type")
	uploadBizType := service.UploadBizType(bizType)
	if uploadBizType == "" {
		uploadBizType = service.BizTypeImage
	}

	// for avatar, use UploadAvatar method with openID naming
	if uploadBizType == service.BizTypeAvatar {
		userID := middleware.GetUserID(c)
		user, err := h.userService.GetByID(userID)
		if err != nil || user == nil {
			c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "用户不存在"))
			return
		}

		avatarURL, err := h.service.UploadAvatar(file, user.OpenID)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, model.Success(gin.H{
			"url": avatarURL,
		}))
		return
	}

	// for trip images, returns public URL
	if uploadBizType == service.BizTypeTrip {
		publicURL, err := h.service.UploadImage(file, uploadBizType)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, model.Success(gin.H{
			"url": publicURL,
		}))
		return
	}

	// for chat images, use UploadImage method
	objectKey, err := h.service.UploadImage(file, uploadBizType)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"key": objectKey,
	}))
}

// GetSignedURL handles GET /api/resource/url
// Query param: key - object key in COS
func (h *UploadHandler) GetSignedURL(c *gin.Context) {
	objectKey := c.Query("key")
	if objectKey == "" {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "key参数不能为空"))
		return
	}

	signedURL, err := h.service.GetSignedURL(objectKey, 3600)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"url": signedURL,
	}))
}
