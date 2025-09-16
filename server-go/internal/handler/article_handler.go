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

type ArticleHandler struct {
	articleService *service.ArticleService
	validator      *validator.Validate
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
		validator:      validator.New(),
	}
}

// GetList 获取文章列表
func (h *ArticleHandler) GetList(c *gin.Context) {
	var req model.ArticleListRequest
	
	// 支持form数据和query参数
	if c.Request.Method == "POST" {
		// 处理POST表单数据
		if tag := c.PostForm("tag"); tag != "" {
			req.Tag = tag
		}
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
		// GET请求使用query参数
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

	articles, err := h.articleService.GetList(c.Request.Context(), &req)
	if err != nil {
		pkg.ServerError(c, "获取文章列表失败")
		return
	}

	pkg.Success(c, articles)
}

// GetByID 获取文章详情
func (h *ArticleHandler) GetByID(c *gin.Context) {
	var idStr string
	
	// 优先从context获取ID（用于兼容接口）
	if id, exists := c.Get("id"); exists {
		idStr = id.(string)
	} else {
		idStr = c.Param("id")
	}
	
	if idStr == "" {
		pkg.ValidateError(c, "缺少文章ID")
		return
	}
	
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		pkg.ValidateError(c, "无效的文章ID")
		return
	}

	article, err := h.articleService.GetByID(c.Request.Context(), id)
	if err != nil {
		pkg.Error(c, "文章不存在")
		return
	}

	pkg.Success(c, article)
}

// GetHot 获取热门文章
func (h *ArticleHandler) GetHot(c *gin.Context) {
	var req model.ArticleHotRequest
	
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
		req.Limit = 8
	}

	if err := h.validator.Struct(&req); err != nil {
		pkg.ValidateError(c, "参数验证失败")
		return
	}

	articles, err := h.articleService.GetHot(c.Request.Context(), &req)
	if err != nil {
		pkg.ServerError(c, "获取热门文章失败")
		return
	}

	pkg.Success(c, articles)
}

// Search 搜索文章
func (h *ArticleHandler) Search(c *gin.Context) {
	var req model.ArticleSearchRequest
	
	// 支持form数据和JSON数据
	if content := c.PostForm("keywords"); content != "" {
		// 表单数据 - keywords字段
		req.Content = content
	} else if content := c.PostForm("content"); content != "" {
		// 表单数据 - content字段
		req.Content = content
	} else if contentFromCtx, exists := c.Get("content"); exists {
		// 从context获取（server.go中设置）
		req.Content = contentFromCtx.(string)
	} else {
		// JSON数据
		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.ValidateError(c, "请求参数错误")
			return
		}
	}

	if err := h.validator.Struct(&req); err != nil {
		pkg.ValidateError(c, "参数验证失败")
		return
	}

	articles, err := h.articleService.Search(c.Request.Context(), &req)
	if err != nil {
		pkg.ServerError(c, "搜索文章失败")
		return
	}

	pkg.Success(c, articles)
}

// GetInfo 获取文章信息统计
func (h *ArticleHandler) GetInfo(c *gin.Context) {
	info, err := h.articleService.GetInfo(c.Request.Context())
	if err != nil {
		pkg.ServerError(c, "获取文章信息失败")
		return
	}

	pkg.Success(c, info)
}

// GetExtend 获取延伸阅读
func (h *ArticleHandler) GetExtend(c *gin.Context) {
	tag := c.PostForm("tag")
	if tag == "" {
		pkg.ValidateError(c, "缺少标签参数")
		return
	}

	// 创建请求结构
	req := &model.ArticleListRequest{
		Tag:   tag,
		Skip:  0,
		Limit: 3, // 延伸阅读一般返回3篇
	}

	articles, err := h.articleService.GetList(c.Request.Context(), req)
	if err != nil {
		pkg.ServerError(c, "获取延伸阅读失败")
		return
	}

	pkg.Success(c, articles)
}