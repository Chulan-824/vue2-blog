package service

import (
	"context"

	"vue2-blog-server/internal/model"
	"vue2-blog-server/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	userRepo    *repository.UserRepository
}

func NewMessageService(messageRepo *repository.MessageRepository, userRepo *repository.UserRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

// Create 创建留言
func (s *MessageService) Create(ctx context.Context, userID primitive.ObjectID, req *model.MessageCreateRequest) error {
	message := &model.Message{
		User:    userID,
		Content: req.Content,
	}

	return s.messageRepo.Create(ctx, message)
}

// AddReply 添加回复
func (s *MessageService) AddReply(ctx context.Context, messageID, userID primitive.ObjectID, req *model.MessageReplyRequest) error {
	child := model.MessageChild{
		User:    userID,
		Content: req.Content,
		ReUser:  req.ReUser,
	}

	return s.messageRepo.AddChildMessage(ctx, messageID, child)
}

// GetList 获取留言列表
func (s *MessageService) GetList(ctx context.Context, req *model.MessageListRequest) ([]*model.MessageDTO, error) {
	messages, err := s.messageRepo.FindList(ctx, req.Skip, req.Limit)
	if err != nil {
		return nil, err
	}

	var dtos []*model.MessageDTO
	for _, message := range messages {
		dto, err := s.toMessageDTO(ctx, message)
		if err != nil {
			continue // 跳过错误的数据
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// toMessageDTO 转换为MessageDTO
func (s *MessageService) toMessageDTO(ctx context.Context, message *model.Message) (*model.MessageDTO, error) {
	// 获取主留言用户信息
	user, err := s.userRepo.FindByID(ctx, message.User)
	if err != nil {
		return nil, err
	}

	dto := &model.MessageDTO{
		ID:      message.ID,
		User:    user.ToDTO(),
		Content: message.Content,
		Date:    message.Date,
	}

	// 处理子留言
	for _, child := range message.Children {
		childUser, err := s.userRepo.FindByID(ctx, child.User)
		if err != nil {
			continue // 跳过找不到用户的子留言
		}

		childDTO := model.MessageChildDTO{
			User:    childUser.ToDTO(),
			Content: child.Content,
			ReUser:  child.ReUser,
			Date:    child.Date,
		}

		dto.Children = append(dto.Children, childDTO)
	}

	return dto, nil
}