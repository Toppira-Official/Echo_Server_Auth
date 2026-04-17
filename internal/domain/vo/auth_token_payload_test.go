package vo

import (
	"testing"
	"time"
)

func TestAuthTokenPayload_NewAuthTokenPayload(t *testing.T) {
	userID, err := NewUserID("123")
	if err != nil {
		t.Fatalf("unexpected userID error: %v", err)
	}

	now := time.Now().UTC()

	tests := []struct {
		name      string
		userID    UserID
		issuedAt  time.Time
		expiredAt time.Time
		expectErr bool
	}{
		{
			name:      "valid payload",
			userID:    userID,
			issuedAt:  now,
			expiredAt: now.Add(1 * time.Hour),
			expectErr: false,
		},
		{
			name:      "invalid payload (expiredAt < issuedAt)",
			userID:    userID,
			issuedAt:  now,
			expiredAt: now.Add(-1 * time.Hour),
			expectErr: true,
		},
		{
			name:      "invalid payload (expiredAt == issuedAt)",
			userID:    userID,
			issuedAt:  now,
			expiredAt: now,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		_, err := NewAuthTokenPayload(tt.userID, tt.issuedAt, tt.expiredAt)

		if tt.expectErr && err == nil {
			t.Fatalf("[%s] expected error but got nil", tt.name)
		}
		if !tt.expectErr && err != nil {
			t.Fatalf("[%s] did NOT expect error but got: %v", tt.name, err)
		}
	}
}

func TestAuthTokenPayload_ExpirationBehavior(t *testing.T) {
	userID, _ := NewUserID("abc")
	issuedAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	expiredAt := issuedAt.Add(2 * time.Hour)

	payload, err := NewAuthTokenPayload(userID, issuedAt, expiredAt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name    string
		checkAt time.Time
		expired bool
		valid   bool
	}{
		{
			name:    "before expiration",
			checkAt: issuedAt.Add(1 * time.Hour),
			expired: false,
			valid:   true,
		},
		{
			name:    "exact expiration moment",
			checkAt: expiredAt,
			expired: false,
			valid:   true,
		},
		{
			name:    "after expiration",
			checkAt: expiredAt.Add(1 * time.Second),
			expired: true,
			valid:   false,
		},
	}

	for _, tt := range tests {
		if payload.IsExpired(tt.checkAt) != tt.expired {
			t.Fatalf("[%s] expected expired=%v but got %v",
				tt.name, tt.expired, payload.IsExpired(tt.checkAt))
		}

		if payload.IsValid(tt.checkAt) != tt.valid {
			t.Fatalf("[%s] expected valid=%v but got %v",
				tt.name, tt.valid, payload.IsValid(tt.checkAt))
		}
	}
}

func TestAuthTokenPayload_Lifetime(t *testing.T) {
	userID, _ := NewUserID("xyz")

	issuedAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	expiredAt := issuedAt.Add(3 * time.Hour)

	payload, _ := NewAuthTokenPayload(userID, issuedAt, expiredAt)

	expectedLifetime := 3 * time.Hour
	if payload.Lifetime() != expectedLifetime {
		t.Fatalf("expected lifetime %v but got %v", expectedLifetime, payload.Lifetime())
	}
}

func TestAuthTokenPayload_Immutability(t *testing.T) {
	userID, _ := NewUserID("immutable")

	issued := time.Now().UTC()
	expired := issued.Add(30 * time.Minute)

	payload, _ := NewAuthTokenPayload(userID, issued, expired)

	retrievedIssued := payload.IssuedAtUTC()
	retrievedIssued = retrievedIssued.Add(9999 * time.Hour)

	if payload.IssuedAtUTC() != issued {
		t.Fatal("AuthTokenPayload is expected to be immutable but state was mutated")
	}
}
