package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Visitor 访客模型
type Visitor struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User primitive.ObjectID `bson:"user" json:"user"`
	Date time.Time          `bson:"date" json:"date"`
}

// VisitorListRequest 访客列表请求
type VisitorListRequest struct {
	Limit int `form:"limit" validate:"min=1,max=50"`
}

// VisitorDTO 访客DTO
type VisitorDTO struct {
	ID   primitive.ObjectID `json:"id"`
	User *UserDTO           `json:"user"`
	Date time.Time          `json:"date"`
}