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

type MessageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *Database) *MessageRepository {
	return &MessageRepository{
		collection: db.GetCollection("message"),
	}
}

// Create 创建留言
func (r *MessageRepository) Create(ctx context.Context, message *model.Message) error {
	message.Date = time.Now()
	message.Children = []model.MessageChild{}
	
	result, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		return err
	}
	
	message.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindList 分页查询留言列表
func (r *MessageRepository) FindList(ctx context.Context, skip, limit int) ([]*model.Message, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "date", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var messages []*model.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	
	return messages, nil
}

// FindByID 根据ID查找留言
func (r *MessageRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Message, error) {
	var message model.Message
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// AddChildMessage 添加子留言
func (r *MessageRepository) AddChildMessage(ctx context.Context, messageID primitive.ObjectID, child model.MessageChild) error {
	child.Date = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": messageID},
		bson.M{"$push": bson.M{"children": child}},
	)
	return err
}