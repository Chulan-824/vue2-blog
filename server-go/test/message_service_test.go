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

type mockMessageRepository struct {
	createFunc   func(ctx context.Context, message *model.Message) error
	addChildFunc func(ctx context.Context, messageID primitive.ObjectID, child model.MessageChild) error
	findListFunc func(ctx context.Context, skip, limit int) ([]*model.Message, error)
}

type mockUserRepository struct {
	findByIDFunc func(ctx context.Context, id primitive.ObjectID) (*model.User, error)
}

var (
	_ service.MessageRepository = (*mockMessageRepository)(nil)
	_ service.UserRepository    = (*mockUserRepository)(nil)
)

func (m *mockMessageRepository) Create(ctx context.Context, message *model.Message) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, message)
	}
	return errors.New("createFunc not implemented")
}

func (m *mockMessageRepository) AddChildMessage(ctx context.Context, messageID primitive.ObjectID, child model.MessageChild) error {
	if m.addChildFunc != nil {
		return m.addChildFunc(ctx, messageID, child)
	}
	return errors.New("addChildFunc not implemented")
}

func (m *mockMessageRepository) FindList(ctx context.Context, skip, limit int) ([]*model.Message, error) {
	if m.findListFunc != nil {
		return m.findListFunc(ctx, skip, limit)
	}
	return nil, errors.New("findListFunc not implemented")
}

func (m *mockUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("findByIDFunc not implemented")
}

func TestMessageServiceCreate(t *testing.T) {
	var captured *model.Message
	messageRepo := &mockMessageRepository{
		createFunc: func(ctx context.Context, message *model.Message) error {
			captured = message
			return nil
		},
	}

	svc := service.NewMessageService(messageRepo, &mockUserRepository{})
	userID := primitive.NewObjectID()
	req := &model.MessageCreateRequest{Content: "Great post!"}
	if err := svc.Create(context.Background(), userID, req); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if captured == nil {
		t.Fatal("expected create to be called")
	}
	if captured.User != userID {
		t.Fatalf("expected message user %s, got %s", userID.Hex(), captured.User.Hex())
	}
	if captured.Content != req.Content {
		t.Fatalf("expected message content %s, got %s", req.Content, captured.Content)
	}
	if captured.Children != nil {
		t.Fatalf("expected children to be nil slice prior to persistence")
	}
}

func TestMessageServiceAddReply(t *testing.T) {
	var (
		capturedMessageID primitive.ObjectID
		capturedChild     model.MessageChild
	)
	messageRepo := &mockMessageRepository{
		addChildFunc: func(ctx context.Context, messageID primitive.ObjectID, child model.MessageChild) error {
			capturedMessageID = messageID
			capturedChild = child
			return nil
		},
	}

	svc := service.NewMessageService(messageRepo, &mockUserRepository{})
	messageID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	req := &model.MessageReplyRequest{Content: "Thanks!", ReUser: "Author"}

	if err := svc.AddReply(context.Background(), messageID, userID, req); err != nil {
		t.Fatalf("AddReply returned error: %v", err)
	}

	if capturedMessageID != messageID {
		t.Fatalf("expected message ID %s, got %s", messageID.Hex(), capturedMessageID.Hex())
	}
	if capturedChild.User != userID {
		t.Fatalf("expected reply user %s, got %s", userID.Hex(), capturedChild.User.Hex())
	}
	if capturedChild.Content != req.Content || capturedChild.ReUser != req.ReUser {
		t.Fatalf("unexpected reply payload: %+v", capturedChild)
	}
}

func TestMessageServiceGetListTransformsMessages(t *testing.T) {
	userID := primitive.NewObjectID()
	childUserID := primitive.NewObjectID()
	missingUserID := primitive.NewObjectID()
	now := time.Now()

	messageRepo := &mockMessageRepository{
		findListFunc: func(ctx context.Context, skip, limit int) ([]*model.Message, error) {
			if skip != 0 || limit != 10 {
				t.Fatalf("unexpected pagination skip=%d limit=%d", skip, limit)
			}
			return []*model.Message{{
				ID:      primitive.NewObjectID(),
				User:    userID,
				Content: "First!",
				Date:    now,
				Children: []model.MessageChild{
					{User: childUserID, Content: "Nice article", ReUser: "tester", Date: now.Add(time.Minute)},
					{User: missingUserID, Content: "spam"},
				},
			}}, nil
		},
	}

	userRepo := &mockUserRepository{
		findByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
			switch id {
			case userID:
				return &model.User{ID: id, User: "owner", Photo: "avatar.jpg"}, nil
			case childUserID:
				return &model.User{ID: id, User: "commenter", Photo: "user.jpg"}, nil
			default:
				return nil, errors.New("not found")
			}
		},
	}

	svc := service.NewMessageService(messageRepo, userRepo)
	messages, err := svc.GetList(context.Background(), &model.MessageListRequest{Limit: 10})
	if err != nil {
		t.Fatalf("GetList returned error: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	msg := messages[0]
	if msg.User == nil || msg.User.User != "owner" {
		t.Fatalf("expected owner user info, got %+v", msg.User)
	}
	if len(msg.Children) != 1 {
		t.Fatalf("expected missing child to be skipped, children=%d", len(msg.Children))
	}
	if msg.Children[0].User == nil || msg.Children[0].User.User != "commenter" {
		t.Fatalf("unexpected child dto: %+v", msg.Children[0])
	}
}
