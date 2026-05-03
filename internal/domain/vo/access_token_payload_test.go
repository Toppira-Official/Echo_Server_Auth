package vo

import (
	"testing"
	"time"
)

func TestAuthTokenPayload_NewAccessTokenPayload(t *testing.T) {
	credentialID, err := NewCredentialID("123")
	if err != nil {
		t.Fatalf("unexpected credentialID error: %v", err)
	}

	now := time.Now().UTC()

	tests := []struct {
		name         string
		credentialID CredentialID
		issuedAt     time.Time
		expiredAt    time.Time
		expectErr    bool
	}{
		{
			name:         "valid payload",
			credentialID: credentialID,
			issuedAt:     now,
			expiredAt:    now.Add(1 * time.Hour),
			expectErr:    false,
		},
		{
			name:         "invalid payload (expiredAt < issuedAt)",
			credentialID: credentialID,
			issuedAt:     now,
			expiredAt:    now.Add(-1 * time.Hour),
			expectErr:    true,
		},
		{
			name:         "invalid payload (expiredAt == issuedAt)",
			credentialID: credentialID,
			issuedAt:     now,
			expiredAt:    now,
			expectErr:    true,
		},
	}

	for _, tt := range tests {
		_, err := NewAccessTokenPayload(tt.credentialID, tt.issuedAt, tt.expiredAt)

		if tt.expectErr && err == nil {
			t.Fatalf("[%s] expected error but got nil", tt.name)
		}
		if !tt.expectErr && err != nil {
			t.Fatalf("[%s] did NOT expect error but got: %v", tt.name, err)
		}
	}
}

func TestAuthTokenPayload_ExpirationBehavior(t *testing.T) {
	credentialID, _ := NewCredentialID("abc")
	issuedAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	expiredAt := issuedAt.Add(2 * time.Hour)

	payload, err := NewAccessTokenPayload(credentialID, issuedAt, expiredAt)
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
	credentialID, _ := NewCredentialID("xyz")

	issuedAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	expiredAt := issuedAt.Add(3 * time.Hour)

	payload, _ := NewAccessTokenPayload(credentialID, issuedAt, expiredAt)

	expectedLifetime := 3 * time.Hour
	if payload.Lifetime() != expectedLifetime {
		t.Fatalf("expected lifetime %v but got %v", expectedLifetime, payload.Lifetime())
	}
}

func TestAuthTokenPayload_Immutability(t *testing.T) {
	credentialID, _ := NewCredentialID("immutable")

	issued := time.Now().UTC()
	expired := issued.Add(30 * time.Minute)

	payload, _ := NewAccessTokenPayload(credentialID, issued, expired)

	if payload.IssuedAtUTC() != issued {
		t.Fatal("AuthTokenPayload is expected to be immutable but state was mutated")
	}
}
