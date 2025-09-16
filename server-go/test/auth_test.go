package test

import (
	"testing"
	"time"

	"vue2-blog-server/internal/auth"
	"vue2-blog-server/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJWTManagerGenerateAndParseToken(t *testing.T) {
	manager := auth.NewJWTManager("secret", time.Minute, time.Hour)
	user := &model.UserDTO{
		ID:    primitive.NewObjectID(),
		User:  "tester",
		Admin: true,
	}

	token, err := manager.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("expected user id %s, got %s", user.ID.Hex(), claims.UserID.Hex())
	}

	if claims.Username != user.User {
		t.Fatalf("expected username %s, got %s", user.User, claims.Username)
	}

	if claims.Admin != user.Admin {
		t.Fatalf("expected admin %t, got %t", user.Admin, claims.Admin)
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		t.Fatalf("expected valid expiration, got %v", claims.ExpiresAt)
	}
}

func TestJWTManagerRefreshToken(t *testing.T) {
	manager := auth.NewJWTManager("secret", time.Minute, time.Hour)
	user := &model.UserDTO{
		ID:    primitive.NewObjectID(),
		User:  "tester",
		Admin: false,
	}

	refreshToken, err := manager.GenerateRefreshToken(user)
	if err != nil {
		t.Fatalf("GenerateRefreshToken returned error: %v", err)
	}

	newToken, err := manager.RefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("RefreshToken returned error: %v", err)
	}

	if newToken == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := manager.ParseToken(newToken)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("expected refreshed token to contain same user id")
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		t.Fatalf("expected refreshed token to have future expiration")
	}
}

func TestJWTManagerParseInvalidToken(t *testing.T) {
	manager := auth.NewJWTManager("secret", time.Minute, time.Hour)
	if _, err := manager.ParseToken("invalid-token"); err == nil {
		t.Fatal("expected error when parsing invalid token")
	}

	if _, err := manager.RefreshToken("invalid-token"); err == nil {
		t.Fatal("expected error when refreshing invalid token")
	}
}
