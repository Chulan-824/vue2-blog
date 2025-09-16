package handler

import (
	"vue2-blog-server/internal/pkg"
	"vue2-blog-server/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadHandler struct {
	uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// UploadAvatar 上传头像
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.AuthError(c, "请先登录")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		pkg.ValidateError(c, "请选择要上传的文件")
		return
	}

	result, err := h.uploadService.UploadAvatar(c.Request.Context(), userID.(primitive.ObjectID), file)
	if err != nil {
		pkg.Error(c, err.Error())
		return
	}

	pkg.SuccessWithMsg(c, "头像上传成功", result)
}

// UploadAvatarLegacy 上传头像（兼容旧版本，使用file字段）
func (h *UploadHandler) UploadAvatarLegacy(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.AuthError(c, "请先登录")
		return
	}

	// 优先尝试 file 字段（前端默认）
	file, err := c.FormFile("file")
	if err != nil {
		// 如果 file 字段不存在，尝试 avatar 字段
		file, err = c.FormFile("avatar")
		if err != nil {
			pkg.ValidateError(c, "请选择要上传的文件")
			return
		}
	}

	result, err := h.uploadService.UploadAvatar(c.Request.Context(), userID.(primitive.ObjectID), file)
	if err != nil {
		pkg.Error(c, err.Error())
		return
	}

	pkg.SuccessWithMsg(c, "头像上传成功", result)
}