package entity

import (
	"auth/internal/domain/vo"
	"testing"
	"time"
)

func TestNewCredential(t *testing.T) {
	now := time.Now().UTC()

	userID, _ := vo.NewUserID("1")

	hash, _ := vo.NewHashedPassword("$2a$10$abcdefghijklmnopqrstuv")

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "valid credential",
			username: "ali",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cred, err := NewCredential(
				userID,
				tt.username,
				now,
				hash,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cred.Username() != tt.username {
				t.Fatalf("expected username %s got %s", tt.username, cred.Username())
			}

			if cred.CreatedAt() != now {
				t.Fatalf("createdAt mismatch")
			}

			if cred.UpdatedAt() != now {
				t.Fatalf("updatedAt mismatch")
			}
		})
	}
}

func TestCredential_ChangePassword(t *testing.T) {

	now := time.Now().UTC()
	newTime := now.Add(time.Hour)

	userID, _ := vo.NewUserID("1")

	hash1, _ := vo.NewHashedPassword("$2a$10$abcdefghijklmnopqrstuv")
	hash2, _ := vo.NewHashedPassword("$2a$10$zzzzzzzzzzzzzzzzzzzzzz")

	cred, _ := NewCredential(
		userID,
		"ali",
		now,
		hash1,
	)

	cred.ChangePassword(hash2, newTime)

	if cred.HashedPassword().Value() != hash2.Value() {
		t.Fatalf("password not updated")
	}

	if cred.UpdatedAt() != newTime {
		t.Fatalf("updatedAt not updated")
	}
}

func TestCredential_UpdateRefreshToken(t *testing.T) {

	now := time.Now().UTC()
	newTime := now.Add(time.Minute)

	userID, _ := vo.NewUserID("1")
	hash, _ := vo.NewHashedPassword("$2a$10$abcdefghijklmnopqrstuv")

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "set refresh token",
			token: "refresh-token-123",
		},
		{
			name:  "replace refresh token",
			token: "new-refresh-token",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			cred, _ := NewCredential(
				userID,
				"ali",
				now,
				hash,
			)

			cred.UpdateRefreshToken(tt.token, newTime)

			if cred.RefreshToken() != tt.token {
				t.Fatalf("expected token %s got %s", tt.token, cred.RefreshToken())
			}

			if cred.UpdatedAt() != newTime {
				t.Fatalf("updatedAt not updated")
			}
		})
	}
}

func TestCredential_RevokeRefreshToken(t *testing.T) {

	now := time.Now().UTC()
	newTime := now.Add(time.Minute)

	userID, _ := vo.NewUserID("1")
	hash, _ := vo.NewHashedPassword("$2a$10$abcdefghijklmnopqrstuv")

	tests := []struct {
		name         string
		initialToken string
		expectEmpty  bool
	}{
		{
			name:         "revoke existing token",
			initialToken: "token123",
			expectEmpty:  true,
		},
		{
			name:         "revoke when no token",
			initialToken: "",
			expectEmpty:  true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			cred, _ := NewCredential(
				userID,
				"ali",
				now,
				hash,
			)

			if tt.initialToken != "" {
				cred.UpdateRefreshToken(tt.initialToken, now)
			}

			cred.RevokeRefreshToken(newTime)

			if cred.RefreshToken() != "" && tt.expectEmpty {
				t.Fatalf("token should be empty")
			}
		})
	}
}
