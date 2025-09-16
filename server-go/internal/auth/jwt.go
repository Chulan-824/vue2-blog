package auth

import (
	"errors"
	"time"

	"vue2-blog-server/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Claims JWT载荷
type Claims struct {
	UserID   primitive.ObjectID `json:"userId"`
	Username string             `json:"username"`
	Admin    bool               `json:"admin"`
	jwt.RegisteredClaims
}

// JWTManager JWT管理器
type JWTManager struct {
	secretKey     string
	tokenExpires  time.Duration
	refreshExpires time.Duration
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secretKey string, tokenExpires, refreshExpires time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenExpires:  tokenExpires,
		refreshExpires: refreshExpires,
	}
}

// GenerateToken 生成访问令牌
func (j *JWTManager) GenerateToken(user *model.UserDTO) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.User,
		Admin:    user.Admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// GenerateRefreshToken 生成刷新令牌
func (j *JWTManager) GenerateRefreshToken(user *model.UserDTO) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.User,
		Admin:    user.Admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ParseToken 解析令牌
func (j *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// RefreshToken 刷新令牌
func (j *JWTManager) RefreshToken(refreshToken string) (string, error) {
	claims, err := j.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	// 创建新的访问令牌
	newClaims := &Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Admin:    claims.Admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExpires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString([]byte(j.secretKey))
}