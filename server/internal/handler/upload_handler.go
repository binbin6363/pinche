package handler

import (
	"net/http"

	"pinche/internal/middleware"
	"pinche/internal/model"
	"pinche/internal/service"

	"github.com/gin-gonic/gin"
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

// Upload handles POST /api/upload
// Query param: biz_type - "images" for chat images, "avatar" for user avatar, "trip" for trip images, "voices" for voice messages
// Response:
//   - for avatar/trip: returns { "url": "public_url" }
//   - for images/voices: returns { "key": "object_key" }
func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "请选择文件"))
		return
	}

	bizType := service.UploadBizType(c.Query("biz_type"))
	if !service.IsValidBizType(bizType) {
		bizType = service.BizTypeImage
	}

	// avatar: use UploadAvatar method with openID naming
	if bizType == service.BizTypeAvatar {
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

	// trip: returns public URL
	if bizType == service.BizTypeTrip {
		publicURL, err := h.service.Upload(file, bizType)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, model.Success(gin.H{
			"url": publicURL,
		}))
		return
	}

	// images/voices: returns object key
	objectKey, err := h.service.Upload(file, bizType)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{
		"key": objectKey,
	}))
}

// UploadImage handles POST /api/upload/image (deprecated, use Upload instead)
func (h *UploadHandler) UploadImage(c *gin.Context) {
	h.Upload(c)
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
