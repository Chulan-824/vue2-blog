package service

import (
	"context"

	"vue2-blog-server/internal/model"
	"vue2-blog-server/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ArticleRepository 定义文章仓储需要实现的方法
type ArticleRepository interface {
	FindList(ctx context.Context, tag string, skip, limit int) ([]*model.Article, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*model.Article, error)
	IncrementPV(ctx context.Context, id primitive.ObjectID) error
	FindHot(ctx context.Context, limit int) ([]*model.Article, error)
	Search(ctx context.Context, content string) ([]*model.Article, error)
	GetInfo(ctx context.Context) (*model.ArticleInfo, error)
}

// 确保实际仓储实现接口
var _ ArticleRepository = (*repository.ArticleRepository)(nil)

type ArticleService struct {
	articleRepo ArticleRepository
}

func NewArticleService(articleRepo ArticleRepository) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
	}
}

// GetList 获取文章列表
func (s *ArticleService) GetList(ctx context.Context, req *model.ArticleListRequest) ([]*model.ArticleDTO, error) {
	articles, err := s.articleRepo.FindList(ctx, req.Tag, req.Skip, req.Limit)
	if err != nil {
		return nil, err
	}

	var dtos []*model.ArticleDTO
	for _, article := range articles {
		dtos = append(dtos, article.ToDTO(false)) // 列表不包含内容
	}

	return dtos, nil
}

// GetByID 根据ID获取文章详情
func (s *ArticleService) GetByID(ctx context.Context, id primitive.ObjectID) (*model.ArticleDTO, error) {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 增加浏览量
	go s.articleRepo.IncrementPV(context.Background(), id)

	return article.ToDTO(true), nil // 详情包含内容
}

// GetHot 获取热门文章
func (s *ArticleService) GetHot(ctx context.Context, req *model.ArticleHotRequest) ([]*model.ArticleDTO, error) {
	if req.Limit == 0 {
		req.Limit = 8 // 默认8篇
	}

	articles, err := s.articleRepo.FindHot(ctx, req.Limit)
	if err != nil {
		return nil, err
	}

	var dtos []*model.ArticleDTO
	for _, article := range articles {
		dtos = append(dtos, article.ToDTO(false))
	}

	return dtos, nil
}

// Search 搜索文章
func (s *ArticleService) Search(ctx context.Context, req *model.ArticleSearchRequest) ([]*model.ArticleDTO, error) {
	articles, err := s.articleRepo.Search(ctx, req.Content)
	if err != nil {
		return nil, err
	}

	var dtos []*model.ArticleDTO
	for _, article := range articles {
		dtos = append(dtos, article.ToDTO(false))
	}

	return dtos, nil
}

// GetInfo 获取文章信息统计
func (s *ArticleService) GetInfo(ctx context.Context) (*model.ArticleInfo, error) {
	return s.articleRepo.GetInfo(ctx)
}
