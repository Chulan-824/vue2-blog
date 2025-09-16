package repository

import (
	"context"
	"time"

	"vue2-blog-server/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{
		collection: db.GetCollection("user"),
	}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	user.RegDate = time.Now().UnixNano() / int64(time.Millisecond)
	if user.Photo == "" {
		user.Photo = "http://47.96.127.142:80/img/defaultPhoto.jpeg"
	}
	
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"user": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePhoto 更新用户头像
func (r *UserRepository) UpdatePhoto(ctx context.Context, userID primitive.ObjectID, photo string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"photo": photo}},
	)
	return err
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"user": username})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}