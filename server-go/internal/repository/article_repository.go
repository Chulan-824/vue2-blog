package repository

import (
	"context"
	"time"

	"vue2-blog-server/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleRepository struct {
	collection     *mongo.Collection
	infoCollection *mongo.Collection
}

func NewArticleRepository(db *Database) *ArticleRepository {
	return &ArticleRepository{
		collection:     db.GetCollection("article"),
		infoCollection: db.GetCollection("articleInfo"),
	}
}

// Create 创建文章
func (r *ArticleRepository) Create(ctx context.Context, article *model.Article) error {
	now := time.Now()
	article.Date = now
	article.UpdateDate = now
	article.PV = 0
	
	if article.Surface == "" {
		article.Surface = "http://47.96.127.142:80/img/defaultSurface.jpg"
	}
	
	result, err := r.collection.InsertOne(ctx, article)
	if err != nil {
		return err
	}
	
	article.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID 根据ID查找文章
func (r *ArticleRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Article, error) {
	var article model.Article
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// FindList 分页查询文章列表
func (r *ArticleRepository) FindList(ctx context.Context, tag string, skip, limit int) ([]*model.Article, error) {
	filter := bson.M{}
	if tag != "" {
		filter["tag"] = tag
	}
	
	opts := options.Find().
		SetSort(bson.D{{Key: "date", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var articles []*model.Article
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	
	return articles, nil
}

// FindHot 查找热门文章
func (r *ArticleRepository) FindHot(ctx context.Context, limit int) ([]*model.Article, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "pv", Value: -1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var articles []*model.Article
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	
	return articles, nil
}

// Search 搜索文章
func (r *ArticleRepository) Search(ctx context.Context, content string) ([]*model.Article, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": content, "$options": "i"}},
			{"tag": bson.M{"$regex": content, "$options": "i"}},
		},
	}
	
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var articles []*model.Article
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	
	return articles, nil
}

// IncrementPV 增加文章浏览量
func (r *ArticleRepository) IncrementPV(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$inc": bson.M{"pv": 1}},
	)
	return err
}

// GetInfo 获取文章信息统计
func (r *ArticleRepository) GetInfo(ctx context.Context) (*model.ArticleInfo, error) {
	var info model.ArticleInfo
	err := r.infoCollection.FindOne(ctx, bson.M{}).Decode(&info)
	if err != nil {
		// 如果没有找到，返回默认值
		if err == mongo.ErrNoDocuments {
			count, err := r.collection.CountDocuments(ctx, bson.M{})
			if err != nil {
				return nil, err
			}
			
			return &model.ArticleInfo{
				Tags: []string{"HTML&Css", "JavaScript", "Node", "Vue&React", "Other"},
				Num:  count,
			}, nil
		}
		return nil, err
	}
	return &info, nil
}

// UpdateInfo 更新文章信息统计
func (r *ArticleRepository) UpdateInfo(ctx context.Context, info *model.ArticleInfo) error {
	opts := options.Replace().SetUpsert(true)
	_, err := r.infoCollection.ReplaceOne(ctx, bson.M{}, info, opts)
	return err
}