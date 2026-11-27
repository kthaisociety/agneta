package signin

import (
	"sync"
	"testing"
	"time"
)

// t0 is a fixed "now" so tests can move time forward without waiting.
var t0 = time.Date(2026, 11, 26, 12, 0, 0, 0, time.UTC)

// signedIn returns a store with one grant, its ID, and a fresh pair of tokens issued at t0.
func signedIn(t *testing.T) (s *Store, grantID, access, refresh string) {
	t.Helper()
	s = NewStore()
	grantID = s.AddGrant(Grant{Subject: "is-user-1", Email: "ada@example.com", RefreshToken: "is-refresh"})
	if grantID == "" {
		t.Fatal("AddGrant returned an empty ID")
	}
	access, refresh, err := s.IssueTokens(grantID, t0)
	if err != nil {
		t.Fatalf("IssueTokens: %v", err)
	}
	if access == "" || refresh == "" || access == refresh {
		t.Fatalf("IssueTokens gave access %q and refresh %q, want two different non-empty tokens", access, refresh)
	}
	return s, grantID, access, refresh
}

// TestAccessTokenFindsItsGrant checks the normal case: a fresh access token leads to its grant.
func TestAccessTokenFindsItsGrant(t *testing.T) {
	s, grantID, access, _ := signedIn(t)
	got, ok := s.LookupAccess(access, t0.Add(time.Minute))
	if !ok || got != grantID {
		t.Errorf("LookupAccess = %q, %v; want %q, true", got, ok, grantID)
	}
}

// TestAccessTokenRejections checks every way an access token must fail.
func TestAccessTokenRejections(t *testing.T) {
	s, _, access, refresh := signedIn(t)
	tests := []struct {
		name  string
		token string
		now   time.Time
	}{
		{"unknown", "not-a-token", t0},
		{"expired", access, t0.Add(accessLifetime + time.Second)},
		{"refresh token used as access", refresh, t0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := s.LookupAccess(tt.token, tt.now); ok {
				t.Errorf("LookupAccess accepted it, grant %q", got)
			}
		})
	}
}

// TestRefreshTokenWorksOnce checks that a refresh token finds its grant once, then never again.
func TestRefreshTokenWorksOnce(t *testing.T) {
	s, grantID, _, refresh := signedIn(t)
	if got, ok := s.TakeRefresh(refresh, t0.Add(time.Hour)); !ok || got != grantID {
		t.Fatalf("first TakeRefresh = %q, %v; want %q, true", got, ok, grantID)
	}
	if _, ok := s.TakeRefresh(refresh, t0.Add(time.Hour)); ok {
		t.Error("second TakeRefresh accepted the same token")
	}
}

// TestRefreshTokenExpires checks that a refresh token stops working after refreshLifetime.
func TestRefreshTokenExpires(t *testing.T) {
	s, _, _, refresh := signedIn(t)
	if _, ok := s.TakeRefresh(refresh, t0.Add(refreshLifetime+time.Second)); ok {
		t.Error("TakeRefresh accepted an expired token")
	}
}

// TestDeleteGrantEndsItsTokens checks that deleting a grant makes both its tokens useless.
func TestDeleteGrantEndsItsTokens(t *testing.T) {
	s, grantID, access, refresh := signedIn(t)
	s.DeleteGrant(grantID)
	if _, ok := s.LookupAccess(access, t0); ok {
		t.Error("LookupAccess accepted a token of a deleted grant")
	}
	if _, ok := s.TakeRefresh(refresh, t0); ok {
		t.Error("TakeRefresh accepted a token of a deleted grant")
	}
}

// TestIssueTokensNeedsAGrant checks that tokens can't be issued for a grant that doesn't exist.
func TestIssueTokensNeedsAGrant(t *testing.T) {
	if _, _, err := NewStore().IssueTokens("no-such-grant", t0); err != errNoGrant {
		t.Errorf("err = %v, want errNoGrant", err)
	}
}

// TestTokensAreStoredHashed checks that the raw tokens never appear as keys in the store.
func TestTokensAreStoredHashed(t *testing.T) {
	s, _, access, refresh := signedIn(t)
	if _, raw := s.access[access]; raw {
		t.Error("access token stored as given, want only its hash")
	}
	if _, raw := s.refresh[refresh]; raw {
		t.Error("refresh token stored as given, want only its hash")
	}
	if _, hashed := s.access[hashToken(access)]; !hashed {
		t.Error("access token not stored under its hash")
	}
}

// TestStoreIsSafeForConcurrentUse runs many sign-ins at once. Without the mutex, `go test -race`
// (which `make test` uses) reports a data race here, or the runtime crashes on a concurrent map write.
func TestStoreIsSafeForConcurrentUse(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := s.AddGrant(Grant{Subject: "is-user", Email: "ada@example.com"})
			access, refresh, err := s.IssueTokens(id, t0)
			if err != nil {
				t.Errorf("IssueTokens: %v", err)
				return
			}
			s.LookupAccess(access, t0)
			s.TakeRefresh(refresh, t0)
			s.DeleteGrant(id)
		}()
	}
	wg.Wait()
}
