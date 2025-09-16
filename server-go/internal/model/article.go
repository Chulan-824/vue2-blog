package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Article 文章模型
type Article struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Type       string               `bson:"type" json:"type" validate:"required,oneof=原创 转载"`
	Title      string               `bson:"title" json:"title" validate:"required"`
	Content    string               `bson:"content" json:"content" validate:"required"`
	Tag        string               `bson:"tag" json:"tag" validate:"required"`
	UpdateDate time.Time            `bson:"updateDate" json:"updateDate"`
	Date       time.Time            `bson:"date" json:"date"`
	Surface    string               `bson:"surface" json:"surface"`
	PV         int64                `bson:"pv" json:"pv"`
	CommentIDs []primitive.ObjectID `bson:"comment" json:"comment"`
}

// ArticleInfo 文章信息统计
type ArticleInfo struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Tags []string           `bson:"tags" json:"tags"`
	Num  int64              `bson:"num" json:"num"`
}

// ArticleListRequest 文章列表请求
type ArticleListRequest struct {
	Tag   string `form:"tag"`
	Skip  int    `form:"skip" validate:"min=0"`
	Limit int    `form:"limit" validate:"min=1,max=20"`
}

// ArticleSearchRequest 文章搜索请求
type ArticleSearchRequest struct {
	Content string `json:"content" validate:"required"`
}

// ArticleHotRequest 热门文章请求
type ArticleHotRequest struct {
	Limit int `form:"limit" validate:"min=1,max=20"`
}

// ArticleDTO 文章传输对象
type ArticleDTO struct {
	ID         primitive.ObjectID `json:"id"`
	Type       string             `json:"type"`
	Title      string             `json:"title"`
	Content    string             `json:"content,omitempty"` // 列表时不返回内容
	Tag        string             `json:"tag"`
	UpdateDate time.Time          `json:"updateDate"`
	Date       time.Time          `json:"date"`
	Surface    string             `json:"surface"`
	PV         int64              `json:"pv"`
}

// ToDTO 转换为DTO
func (a *Article) ToDTO(includeContent bool) *ArticleDTO {
	dto := &ArticleDTO{
		ID:         a.ID,
		Type:       a.Type,
		Title:      a.Title,
		Tag:        a.Tag,
		UpdateDate: a.UpdateDate,
		Date:       a.Date,
		Surface:    a.Surface,
		PV:         a.PV,
	}
	
	if includeContent {
		dto.Content = a.Content
	}
	
	return dto
}