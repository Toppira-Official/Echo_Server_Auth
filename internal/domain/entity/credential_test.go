package entity

import (
	"auth/internal/domain/vo"
	"testing"
	"time"
)

func TestNewCredential(t *testing.T) {
	now := time.Now().UTC()

	credentialID, _ := vo.NewCredentialID("1")

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
				credentialID,
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

	credentialID, _ := vo.NewCredentialID("1")

	hash1, _ := vo.NewHashedPassword("$2a$10$abcdefghijklmnopqrstuv")
	hash2, _ := vo.NewHashedPassword("$2a$10$zzzzzzzzzzzzzzzzzzzzzz")

	cred, _ := NewCredential(
		credentialID,
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
