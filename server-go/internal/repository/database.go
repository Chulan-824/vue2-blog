package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database 数据库连接结构
type Database struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewDatabase 创建数据库连接
func NewDatabase(uri, dbName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Database{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}

// GetCollection 获取集合
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.db.Collection(name)
}