package handler

import (
	"strconv"

	"vue2-blog-server/internal/model"
	"vue2-blog-server/internal/pkg"
	"vue2-blog-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageHandler struct {
	messageService *service.MessageService
	validator      *validator.Validate
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		validator:      validator.New(),
	}
}

// Create 创建留言
func (h *MessageHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.AuthError(c, "请先登录")
		return
	}

	var req model.MessageCreateRequest
	
	// 支持form数据和JSON数据
	if content := c.PostForm("content"); content != "" {
		req.Content = content
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.ValidateError(c, "请求参数错误")
			return
		}
	}

	if err := h.validator.Struct(&req); err != nil {
		pkg.ValidateError(c, "参数验证失败")
		return
	}

	err := h.messageService.Create(c.Request.Context(), userID.(primitive.ObjectID), &req)
	if err != nil {
		pkg.ServerError(c, "创建留言失败")
		return
	}

	pkg.SuccessWithMsg(c, "留言成功", nil)
}

// Reply 回复留言
func (h *MessageHandler) Reply(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.AuthError(c, "请先登录")
		return
	}

	var messageIDStr string
	// 优先从context获取ID（用于兼容接口）
	if id, exists := c.Get("id"); exists {
		messageIDStr = id.(string)
	} else {
		messageIDStr = c.Param("id")
	}
	
	if messageIDStr == "" {
		pkg.ValidateError(c, "缺少留言ID")
		return
	}
	
	messageID, err := primitive.ObjectIDFromHex(messageIDStr)
	if err != nil {
		pkg.ValidateError(c, "无效的留言ID")
		return
	}

	var req model.MessageReplyRequest
	
	// 支持form数据和JSON数据
	if content := c.PostForm("content"); content != "" {
		req.Content = content
		req.ReUser = c.PostForm("reUser") // 可选字段
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.ValidateError(c, "请求参数错误")
			return
		}
	}

	if err := h.validator.Struct(&req); err != nil {
		pkg.ValidateError(c, "参数验证失败")
		return
	}

	err = h.messageService.AddReply(c.Request.Context(), messageID, userID.(primitive.ObjectID), &req)
	if err != nil {
		pkg.ServerError(c, "回复留言失败")
		return
	}

	pkg.SuccessWithMsg(c, "回复成功", nil)
}

// GetList 获取留言列表
func (h *MessageHandler) GetList(c *gin.Context) {
	var req model.MessageListRequest
	
	// 支持form数据和query参数
	if c.Request.Method == "POST" {
		if skip := c.PostForm("skip"); skip != "" {
			if s, err := strconv.Atoi(skip); err == nil {
				req.Skip = s
			}
		}
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
		req.Limit = 5
	}

	if err := h.validator.Struct(&req); err != nil {
		pkg.ValidateError(c, "参数验证失败")
		return
	}

	messages, err := h.messageService.GetList(c.Request.Context(), &req)
	if err != nil {
		pkg.ServerError(c, "获取留言列表失败")
		return
	}

	pkg.Success(c, messages)
}