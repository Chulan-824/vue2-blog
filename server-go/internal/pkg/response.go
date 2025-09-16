package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 响应码定义
const (
	SuccessCode = 0
	ErrorCode   = 1
	ValidateCode = 2
	AuthCode    = 3
	ServerCode  = 4
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SuccessCode,
		Msg:  "操作成功",
		Data: data,
	})
}

// SuccessWithMsg 带消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SuccessCode,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: ErrorCode,
		Msg:  msg,
	})
}

// ValidateError 验证错误响应
func ValidateError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: ValidateCode,
		Msg:  msg,
	})
}

// AuthError 认证错误响应
func AuthError(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: AuthCode,
		Msg:  msg,
	})
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: ServerCode,
		Msg:  msg,
	})
}