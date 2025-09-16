package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User     string             `bson:"user" json:"user" validate:"required,min=2,max=7"`
	Pwd      string             `bson:"pwd" json:"pwd" validate:"required,min=6,max=18"`
	RegDate  int64              `bson:"regDate" json:"regDate"`
	Photo    string             `bson:"photo" json:"photo"`
	Disabled bool               `bson:"disabled" json:"disabled"`
	Admin    bool               `bson:"admin" json:"admin"`
}

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID       primitive.ObjectID `json:"id"`
	User     string             `json:"user"`
	RegDate  int64              `json:"regDate"`
	Photo    string             `json:"photo"`
	Disabled bool               `json:"disabled"`
	Admin    bool               `json:"admin"`
}

// UserLoginRequest 登录请求
type UserLoginRequest struct {
	User string `json:"user" validate:"required,min=2,max=7"`
	Pwd  string `json:"pwd" validate:"required,min=6,max=18"`
}

// UserRegisterRequest 注册请求
type UserRegisterRequest struct {
	User  string `json:"user" validate:"required,min=2,max=7"`
	Pwd   string `json:"pwd" validate:"required,min=6,max=18"`
	Vcode string `json:"vcode" validate:"required"`
}

// ToDTO 转换为DTO
func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:       u.ID,
		User:     u.User,
		RegDate:  u.RegDate,
		Photo:    u.Photo,
		Disabled: u.Disabled,
		Admin:    u.Admin,
	}
}