package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"vue2-blog-server/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadService struct {
	userService *UserService
	uploadPath  string
	baseURL     string
}

func NewUploadService(userService *UserService, uploadPath, baseURL string) *UploadService {
	return &UploadService{
		userService: userService,
		uploadPath:  uploadPath,
		baseURL:     baseURL,
	}
}

// UploadAvatar 上传头像
func (s *UploadService) UploadAvatar(ctx context.Context, userID primitive.ObjectID, file *multipart.FileHeader) (*model.UploadAvatarResponse, error) {
	// 检查文件类型
	if !s.isValidImageType(file.Filename) {
		return nil, fmt.Errorf("不支持的文件类型")
	}

	// 检查文件大小 (5MB)
	if file.Size > 5*1024*1024 {
		return nil, fmt.Errorf("文件大小不能超过5MB")
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", userID.Hex(), time.Now().UnixNano(), ext)

	// 确保上传目录存在
	avatarDir := filepath.Join(s.uploadPath, "avatar")
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %v", err)
	}

	// 保存文件
	dst := filepath.Join(avatarDir, filename)
	if err := s.saveFile(file, dst); err != nil {
		return nil, fmt.Errorf("保存文件失败: %v", err)
	}

	// 生成访问URL - 使用/avatar路径，对应server.go中的静态路由映射
	avatarURL := fmt.Sprintf("%s/avatar/%s", s.baseURL, filename)

	// 更新用户头像
	if err := s.userService.UpdateAvatar(ctx, userID, avatarURL); err != nil {
		// 删除已上传的文件
		os.Remove(dst)
		return nil, fmt.Errorf("更新用户头像失败: %v", err)
	}

	return &model.UploadAvatarResponse{
		URL: avatarURL,
	}, nil
}

// isValidImageType 检查是否为有效的图片类型
func (s *UploadService) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// saveFile 保存文件
func (s *UploadService) saveFile(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}