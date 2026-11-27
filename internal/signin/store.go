// Package signin lets Advisors sign in: agneta is Claude's OAuth authorization server and a client of
// identity-service (docs/adr/0003).
package signin

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// How long agneta's tokens stay valid. Claude refreshes the access token about once an hour.
const (
	accessLifetime  = time.Hour
	refreshLifetime = 30 * 24 * time.Hour
)

// errNoGrant is returned when tokens are asked for a grant the store doesn't have.
var errNoGrant = errors.New("no such grant")

// Grant is one Advisor's sign-in, who they are at identity-service.
type Grant struct {
	// Subject: Advisor's user ID at identity-service.
	Subject string
	Email   string
	// RefreshToken is the identity-service refresh token, used each time Claude refreshes (step 9).
	RefreshToken string
}

// token is what a stored access or refresh token points to.
type token struct {
	grantID string
	expires time.Time
}

// Store holds grants and agneta's tokens. Tokens are kept only as hashes (hashToken), never as given out.
// Each HTTP request runs in its own goroutine, and they all share one Store.
// In memory, a restart disconnects every Advisor; move to Postgres together with report storage.
type Store struct {
	mu      sync.Mutex
	grants  map[string]Grant // grant ID → grant
	access  map[string]token // hashToken(access token) → token
	refresh map[string]token // hashToken(refresh token) → token
}

// NewStore returns an empty store, ready to use.
func NewStore() *Store {
	return &Store{
		grants:  make(map[string]Grant),
		access:  make(map[string]token),
		refresh: make(map[string]token),
	}
}

// AddGrant stores g under a new random ID and returns that ID.
func (s *Store) AddGrant(g Grant) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	grantID := rand.Text()
	s.grants[grantID] = g
	return grantID
}

// IssueTokens gives out a new access and refresh token for the grant, valid from now.
// It returns errNoGrant if the grant doesn't exist.
func (s *Store) IssueTokens(grantID string, now time.Time) (access, refresh string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.grants[grantID]; !ok {
		return "", "", errNoGrant
	}

	access = rand.Text()
	refresh = rand.Text()

	s.access[hashToken(access)] = token{grantID: grantID, expires: now.Add(accessLifetime)}
	s.refresh[hashToken(refresh)] = token{grantID: grantID, expires: now.Add(refreshLifetime)}
	return access, refresh, nil
}

// LookupAccess returns the grant ID behind an access token, if the token is known, not expired at now,
// and its grant still exists.
func (s *Store) LookupAccess(access string, now time.Time) (grantID string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	accessToken, found := s.access[hashToken(access)]
	if !found {
		return "", false
	}
	if now.After(accessToken.expires) {
		return "", false
	}
	if _, found := s.grants[accessToken.grantID]; !found {
		return "", false
	}

	return accessToken.grantID, true
}

// TakeRefresh returns the grant ID behind a refresh token, under the same conditions as LookupAccess,
// and removes the token: a refresh token works only once.
func (s *Store) TakeRefresh(refresh string, now time.Time) (grantID string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := hashToken(refresh)
	refreshToken, found := s.refresh[key]
	if !found {
		return "", false
	}
	// Remove it before any other check: a refresh token is spent once it's presented, valid or not.
	delete(s.refresh, key)

	if now.After(refreshToken.expires) {
		return "", false
	}
	if _, found := s.grants[refreshToken.grantID]; !found {
		return "", false
	}
	return refreshToken.grantID, true
}

// DeleteGrant removes a grant. Its tokens stop working because the lookups check that the grant exists.
func (s *Store) DeleteGrant(grantID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.grants, grantID)
}

// hashToken returns the SHA-256 of a raw token as hex: the key it's stored under.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
