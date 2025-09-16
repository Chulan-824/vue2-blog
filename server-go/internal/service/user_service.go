package service

import (
	"context"
	"errors"

	"vue2-blog-server/internal/model"
	"vue2-blog-server/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    *repository.UserRepository
	visitorRepo *repository.VisitorRepository
}

func NewUserService(userRepo *repository.UserRepository, visitorRepo *repository.VisitorRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
		visitorRepo: visitorRepo,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, req *model.UserRegisterRequest) error {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(ctx, req.User)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		User:     req.User,
		Pwd:      string(hashedPassword),
		Disabled: false,
		Admin:    false,
	}

	return s.userRepo.Create(ctx, user)
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.UserDTO, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.User)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 检查密码（兼容明文密码）
	err = s.checkPassword(req.Pwd, user.Pwd)
	if err != nil {
		return nil, errors.New("密码错误")
	}

	// 如果是明文密码，更新为加密密码
	if s.isPlaintextPassword(user.Pwd) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Pwd), bcrypt.DefaultCost)
		if err == nil {
			user.Pwd = string(hashedPassword)
			// 这里可以异步更新密码，避免影响登录响应时间
		}
	}

	// 记录访客
	go s.recordVisitor(context.Background(), user.ID)

	return user.ToDTO(), nil
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(ctx context.Context, userID primitive.ObjectID) (*model.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.ToDTO(), nil
}

// UpdateAvatar 更新用户头像
func (s *UserService) UpdateAvatar(ctx context.Context, userID primitive.ObjectID, avatarURL string) error {
	return s.userRepo.UpdatePhoto(ctx, userID, avatarURL)
}

// checkPassword 检查密码
func (s *UserService) checkPassword(inputPassword, storedPassword string) error {
	// 如果是明文密码（兼容旧数据）
	if s.isPlaintextPassword(storedPassword) {
		if inputPassword == storedPassword {
			return nil
		}
		return errors.New("密码错误")
	}

	// bcrypt加密密码
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
}

// isPlaintextPassword 判断是否为明文密码
func (s *UserService) isPlaintextPassword(password string) bool {
	// bcrypt hash通常以$2a$, $2b$, $2y$ 开头
	return len(password) < 60 || password[:3] != "$2a" && password[:3] != "$2b" && password[:3] != "$2y"
}

// recordVisitor 记录访客
func (s *UserService) recordVisitor(ctx context.Context, userID primitive.ObjectID) {
	// 先删除该用户的旧访客记录
	s.visitorRepo.DeleteByUserID(ctx, userID)

	// 添加新的访客记录
	visitor := &model.Visitor{
		User: userID,
	}
	s.visitorRepo.Create(ctx, visitor)
}