package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageChild 子留言
type MessageChild struct {
	User    primitive.ObjectID `bson:"user" json:"user"`
	Content string             `bson:"content" json:"content"`
	ReUser  string             `bson:"reUser" json:"reUser"`
	Date    time.Time          `bson:"date" json:"date"`
}

// Message 留言模型
type Message struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User     primitive.ObjectID `bson:"user" json:"user"`
	Content  string             `bson:"content" json:"content" validate:"required"`
	Date     time.Time          `bson:"date" json:"date"`
	Children []MessageChild     `bson:"children" json:"children"`
}

// MessageCreateRequest 创建留言请求
type MessageCreateRequest struct {
	Content string `json:"content" validate:"required"`
}

// MessageReplyRequest 回复留言请求
type MessageReplyRequest struct {
	Content string `json:"content" validate:"required"`
	ReUser  string `json:"reUser"`
}

// MessageListRequest 留言列表请求
type MessageListRequest struct {
	Skip  int `form:"skip" validate:"min=0"`
	Limit int `form:"limit" validate:"min=1,max=20"`
}

// MessageChildDTO 子留言DTO
type MessageChildDTO struct {
	User    *UserDTO  `json:"user"`
	Content string    `json:"content"`
	ReUser  string    `json:"reUser"`
	Date    time.Time `json:"date"`
}

// MessageDTO 留言DTO
type MessageDTO struct {
	ID       primitive.ObjectID `json:"id"`
	User     *UserDTO           `json:"user"`
	Content  string             `json:"content"`
	Date     time.Time          `json:"date"`
	Children []MessageChildDTO  `json:"children"`
}