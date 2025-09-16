package service

import (
	"context"

	"vue2-blog-server/internal/model"
	"vue2-blog-server/internal/repository"
)

type VisitorService struct {
	visitorRepo *repository.VisitorRepository
	userRepo    *repository.UserRepository
}

func NewVisitorService(visitorRepo *repository.VisitorRepository, userRepo *repository.UserRepository) *VisitorService {
	return &VisitorService{
		visitorRepo: visitorRepo,
		userRepo:    userRepo,
	}
}

// GetList 获取访客列表
func (s *VisitorService) GetList(ctx context.Context, req *model.VisitorListRequest) ([]*model.VisitorDTO, error) {
	if req.Limit == 0 {
		req.Limit = 12 // 默认12个
	}

	visitors, err := s.visitorRepo.FindList(ctx, req.Limit)
	if err != nil {
		return nil, err
	}

	var dtos []*model.VisitorDTO
	for _, visitor := range visitors {
		dto, err := s.toVisitorDTO(ctx, visitor)
		if err != nil {
			continue // 跳过错误的数据
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

// toVisitorDTO 转换为VisitorDTO
func (s *VisitorService) toVisitorDTO(ctx context.Context, visitor *model.Visitor) (*model.VisitorDTO, error) {
	user, err := s.userRepo.FindByID(ctx, visitor.User)
	if err != nil {
		return nil, err
	}

	return &model.VisitorDTO{
		ID:   visitor.ID,
		User: user.ToDTO(),
		Date: visitor.Date,
	}, nil
}