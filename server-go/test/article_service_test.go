package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"vue2-blog-server/internal/model"
	"vue2-blog-server/internal/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockArticleRepository struct {
	findListFunc    func(ctx context.Context, tag string, skip, limit int) ([]*model.Article, error)
	findByIDFunc    func(ctx context.Context, id primitive.ObjectID) (*model.Article, error)
	incrementPVFunc func(ctx context.Context, id primitive.ObjectID) error
	findHotFunc     func(ctx context.Context, limit int) ([]*model.Article, error)
	searchFunc      func(ctx context.Context, content string) ([]*model.Article, error)
	getInfoFunc     func(ctx context.Context) (*model.ArticleInfo, error)
}

var _ service.ArticleRepository = (*mockArticleRepository)(nil)

func (m *mockArticleRepository) FindList(ctx context.Context, tag string, skip, limit int) ([]*model.Article, error) {
	if m.findListFunc != nil {
		return m.findListFunc(ctx, tag, skip, limit)
	}
	return nil, errors.New("findListFunc not implemented")
}

func (m *mockArticleRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Article, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("findByIDFunc not implemented")
}

func (m *mockArticleRepository) IncrementPV(ctx context.Context, id primitive.ObjectID) error {
	if m.incrementPVFunc != nil {
		return m.incrementPVFunc(ctx, id)
	}
	return nil
}

func (m *mockArticleRepository) FindHot(ctx context.Context, limit int) ([]*model.Article, error) {
	if m.findHotFunc != nil {
		return m.findHotFunc(ctx, limit)
	}
	return nil, errors.New("findHotFunc not implemented")
}

func (m *mockArticleRepository) Search(ctx context.Context, content string) ([]*model.Article, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, content)
	}
	return nil, errors.New("searchFunc not implemented")
}

func (m *mockArticleRepository) GetInfo(ctx context.Context) (*model.ArticleInfo, error) {
	if m.getInfoFunc != nil {
		return m.getInfoFunc(ctx)
	}
	return nil, errors.New("getInfoFunc not implemented")
}

func TestArticleServiceGetListTransformsArticles(t *testing.T) {
	repo := &mockArticleRepository{
		findListFunc: func(ctx context.Context, tag string, skip, limit int) ([]*model.Article, error) {
			if tag != "Go" {
				t.Fatalf("expected tag Go, got %s", tag)
			}
			if skip != 2 || limit != 5 {
				t.Fatalf("unexpected pagination values skip=%d limit=%d", skip, limit)
			}
			return []*model.Article{{
				ID:      primitive.NewObjectID(),
				Type:    "原创",
				Title:   "Concurrency Patterns",
				Content: "goroutines",
				Tag:     tag,
				Date:    time.Now(),
			}}, nil
		},
	}

	service := service.NewArticleService(repo)
	req := &model.ArticleListRequest{Tag: "Go", Skip: 2, Limit: 5}
	articles, err := service.GetList(context.Background(), req)
	if err != nil {
		t.Fatalf("GetList returned error: %v", err)
	}
	if len(articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(articles))
	}
	if articles[0].Content != "" {
		t.Fatalf("expected content to be omitted in list response")
	}
}

func TestArticleServiceGetByIDIncrementsPV(t *testing.T) {
	articleID := primitive.NewObjectID()
	incrementCh := make(chan primitive.ObjectID, 1)
	repo := &mockArticleRepository{
		findByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*model.Article, error) {
			if id != articleID {
				t.Fatalf("expected id %s, got %s", articleID.Hex(), id.Hex())
			}
			return &model.Article{
				ID:      id,
				Type:    "原创",
				Title:   "Understanding Channels",
				Content: "channels are powerful",
				Tag:     "Go",
				Date:    time.Now(),
			}, nil
		},
		incrementPVFunc: func(ctx context.Context, id primitive.ObjectID) error {
			incrementCh <- id
			return nil
		},
	}

	service := service.NewArticleService(repo)
	article, err := service.GetByID(context.Background(), articleID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if article.Content == "" {
		t.Fatalf("expected detailed view to include content")
	}

	select {
	case got := <-incrementCh:
		if got != articleID {
			t.Fatalf("expected increment for %s, got %s", articleID.Hex(), got.Hex())
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected IncrementPV to be called")
	}
}

func TestArticleServiceGetHotUsesDefaultLimit(t *testing.T) {
	capturedLimit := 0
	repo := &mockArticleRepository{
		findHotFunc: func(ctx context.Context, limit int) ([]*model.Article, error) {
			capturedLimit = limit
			return []*model.Article{{
				ID:    primitive.NewObjectID(),
				Type:  "原创",
				Title: "Hot Article",
				Tag:   "Go",
				Date:  time.Now(),
			}}, nil
		},
	}

	service := service.NewArticleService(repo)
	req := &model.ArticleHotRequest{}
	articles, err := service.GetHot(context.Background(), req)
	if err != nil {
		t.Fatalf("GetHot returned error: %v", err)
	}
	if capturedLimit != 8 {
		t.Fatalf("expected default limit 8, got %d", capturedLimit)
	}
	if len(articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(articles))
	}
	if articles[0].Content != "" {
		t.Fatalf("expected content omitted for hot articles")
	}
}

func TestArticleServiceSearchAndGetInfo(t *testing.T) {
	repo := &mockArticleRepository{
		searchFunc: func(ctx context.Context, content string) ([]*model.Article, error) {
			if content != "golang" {
				t.Fatalf("expected search term golang, got %s", content)
			}
			return []*model.Article{{
				ID:      primitive.NewObjectID(),
				Type:    "原创",
				Title:   "Searching Go",
				Content: "golang tips",
				Tag:     "Go",
				Date:    time.Now(),
			}}, nil
		},
		getInfoFunc: func(ctx context.Context) (*model.ArticleInfo, error) {
			return &model.ArticleInfo{Tags: []string{"Go"}, Num: 1}, nil
		},
	}

	service := service.NewArticleService(repo)
	articles, err := service.Search(context.Background(), &model.ArticleSearchRequest{Content: "golang"})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(articles))
	}

	info, err := service.GetInfo(context.Background())
	if err != nil {
		t.Fatalf("GetInfo returned error: %v", err)
	}
	if info.Num != 1 || len(info.Tags) != 1 || info.Tags[0] != "Go" {
		t.Fatalf("unexpected info payload: %+v", info)
	}
}
