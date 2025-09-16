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

type VisitorRepository struct {
	collection *mongo.Collection
}

func NewVisitorRepository(db *Database) *VisitorRepository {
	return &VisitorRepository{
		collection: db.GetCollection("visitor"),
	}
}

// Create 创建访客记录
func (r *VisitorRepository) Create(ctx context.Context, visitor *model.Visitor) error {
	visitor.Date = time.Now()
	
	result, err := r.collection.InsertOne(ctx, visitor)
	if err != nil {
		return err
	}
	
	visitor.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// DeleteByUserID 删除用户的访客记录
func (r *VisitorRepository) DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user": userID})
	return err
}

// FindList 查询访客列表
func (r *VisitorRepository) FindList(ctx context.Context, limit int) ([]*model.Visitor, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "date", Value: -1}}).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var visitors []*model.Visitor
	if err = cursor.All(ctx, &visitors); err != nil {
		return nil, err
	}
	
	return visitors, nil
}