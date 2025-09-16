package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"sync"

	"vue2-blog-server/internal/auth"
	"vue2-blog-server/internal/model"
	"vue2-blog-server/internal/pkg"
	"vue2-blog-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	userService  *service.UserService
	jwtManager   *auth.JWTManager
	validator    *validator.Validate
	vcodeStorage map[string]string // 临时存储验证码，生产环境应使用Redis
	vcodeMu      sync.RWMutex
}

func NewAuthHandler(userService *service.UserService, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		jwtManager:   jwtManager,
		validator:    validator.New(),
		vcodeStorage: make(map[string]string),
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.UserRegisterRequest

	// 支持form数据和JSON数据
	if user := c.PostForm("user"); user != "" {
		req.User = user
		req.Pwd = c.PostForm("pwd")
		req.Vcode = c.PostForm("vcode")
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.ValidateError(c, "请求参数错误")
			return
		}
	}

	if err := h.validator.Struct(&req); err != nil {
		pkg.ValidateError(c, "数据验证失败")
		return
	}

	// 验证验证码
	sessionKey := c.GetHeader("X-Session-Key")
	if sessionKey == "" {
		sessionKey = c.PostForm("sessionKey")
	}

	if sessionKey != "" {
		if storedCode, exists := h.getVCode(sessionKey); exists && storedCode != "" {
			if storedCode != req.Vcode {
				pkg.ValidateError(c, "验证码错误")
				return
			}
			// 验证成功后删除验证码
			h.deleteVCode(sessionKey)
		}
	}

	err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		pkg.Error(c, err.Error())
		return
	}

	pkg.SuccessWithMsg(c, "注册成功", nil)
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.UserLoginRequest

	// 支持form数据和JSON数据
	if user := c.PostForm("user"); user != "" {
		req.User = user
		req.Pwd = c.PostForm("pwd")
	} else {
		if err := c.ShouldBindJSON(&req); err != nil {
			pkg.ValidateError(c, "请求参数错误")
			return
		}
	}

	if err := h.validator.Struct(&req); err != nil {
		pkg.ValidateError(c, "数据验证失败")
		return
	}

	user, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		pkg.Error(c, err.Error())
		return
	}

	// 生成JWT token
	token, err := h.jwtManager.GenerateToken(user)
	if err != nil {
		pkg.ServerError(c, "生成令牌失败")
		return
	}

	// 设置Cookie
	c.SetCookie("token", token, 3600*24*7, "/", "", false, true) // 7天有效期

	pkg.SuccessWithMsg(c, "登录成功", gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 清除Cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	pkg.SuccessWithMsg(c, "退出登录成功", nil)
}

// Profile 获取用户资料
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.AuthError(c, "未找到用户信息")
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID.(primitive.ObjectID))
	if err != nil {
		pkg.Error(c, "获取用户信息失败")
		return
	}

	pkg.Success(c, gin.H{"userInfo": user})
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		pkg.ValidateError(c, "缺少刷新令牌")
		return
	}

	newToken, err := h.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		pkg.AuthError(c, "刷新令牌失败")
		return
	}

	pkg.Success(c, gin.H{"token": newToken})
}

// CheckLogin 检查登录状态（兼容旧接口）
func (h *AuthHandler) CheckLogin(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		pkg.Success(c, gin.H{"userInfo": false})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID.(primitive.ObjectID))
	if err != nil {
		pkg.Success(c, gin.H{"userInfo": false})
		return
	}

	pkg.Success(c, gin.H{"userInfo": user})
}

// CheckLoginLegacy 检查登录状态（旧版本兼容格式）
func (h *AuthHandler) CheckLoginLegacy(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		// 直接返回裸对象，不使用pkg.Success包装
		c.JSON(http.StatusOK, gin.H{"userInfo": false})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"userInfo": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userInfo": user})
}

// GenerateVCode 生成验证码
func (h *AuthHandler) GenerateVCode(c *gin.Context) {
	// 生成随机字符串作为验证码
	code, err := h.generateRandomCode(4)
	if err != nil {
		pkg.ServerError(c, "生成验证码失败")
		return
	}

	// 生成会话key
	sessionKey, err := h.generateSessionKey()
	if err != nil {
		pkg.ServerError(c, "生成会话标识失败")
		return
	}

	// 存储验证码
	h.setVCode(sessionKey, code)

	// 生成SVG图片
	svgImg, err := h.generateCodeImage(code)
	if err != nil {
		pkg.ServerError(c, "生成验证码图片失败")
		return
	}

	c.Header("X-Session-Key", sessionKey)
	pkg.Success(c, gin.H{
		"svgCode":    svgImg,
		"sessionKey": sessionKey,
	})
}

// CheckVCode 验证验证码
func (h *AuthHandler) CheckVCode(c *gin.Context) {
	svgCode := c.PostForm("svgCode")
	sessionKey := c.GetHeader("X-Session-Key")
	if sessionKey == "" {
		sessionKey = c.PostForm("sessionKey")
	}

	if sessionKey == "" {
		pkg.ValidateError(c, "缺少会话标识")
		return
	}

	storedCode, exists := h.getVCode(sessionKey)
	if !exists {
		pkg.ValidateError(c, "验证码已过期")
		return
	}

	if storedCode != svgCode {
		pkg.ValidateError(c, "验证码错误")
		return
	}

	pkg.Success(c, gin.H{"valid": true})
}

// generateRandomCode 生成随机验证码
func (h *AuthHandler) generateRandomCode(length int) (string, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	code := make([]byte, length)
	if _, err := rand.Read(code); err != nil {
		return "", err
	}

	for i := range code {
		code[i] = chars[int(code[i])%len(chars)]
	}

	return string(code), nil
}

// generateSessionKey 生成会话key
func (h *AuthHandler) generateSessionKey() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// generateCodeImage 生成验证码图片
func (h *AuthHandler) generateCodeImage(code string) (string, error) {
	// 创建简单的验证码图片
	width, height := 120, 40
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充白色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// 这里应该绘制验证码文字，为了简化，我们返回base64编码的PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (h *AuthHandler) getVCode(sessionKey string) (string, bool) {
	h.vcodeMu.RLock()
	defer h.vcodeMu.RUnlock()

	code, exists := h.vcodeStorage[sessionKey]
	return code, exists
}

func (h *AuthHandler) setVCode(sessionKey, code string) {
	h.vcodeMu.Lock()
	defer h.vcodeMu.Unlock()

	h.vcodeStorage[sessionKey] = code
}

func (h *AuthHandler) deleteVCode(sessionKey string) {
	h.vcodeMu.Lock()
	defer h.vcodeMu.Unlock()

	delete(h.vcodeStorage, sessionKey)
}
