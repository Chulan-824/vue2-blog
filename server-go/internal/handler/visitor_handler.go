package handler

import (
	"strconv"

	"vue2-blog-server/internal/model"
	"vue2-blog-server/internal/pkg"
	"vue2-blog-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type VisitorHandler struct {
	visitorService *service.VisitorService
	validator      *validator.Validate
}

func NewVisitorHandler(visitorService *service.VisitorService) *VisitorHandler {
	return &VisitorHandler{
		visitorService: visitorService,
		validator:      validator.New(),
	}
}

// GetList 获取访客列表
func (h *VisitorHandler) GetList(c *gin.Context) {
	var req model.VisitorListRequest
	
	// 支持form数据和query参数
	if c.Request.Method == "POST" {
		if limit := c.PostForm("limit"); limit != "" {
			if l, err := strconv.Atoi(limit); err == nil {
				req.Limit = l
			}
		}
	} else {
		if err := c.ShouldBindQuery(&req); err != nil {
			pkg.ValidateError(c, "请求参数错误")
			return
		}
	}

	// 设置默认值
	if req.Limit == 0 {
		req.Limit = 12
	}

	if err := h.validator.Struct(&req); err != nil {
		pkg.ValidateError(c, "参数验证失败")
		return
	}

	visitors, err := h.visitorService.GetList(c.Request.Context(), &req)
	if err != nil {
		pkg.ServerError(c, "获取访客列表失败")
		return
	}

	pkg.Success(c, visitors)
}