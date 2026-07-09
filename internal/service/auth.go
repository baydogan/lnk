package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
)

type AuthService struct{ keys APIKeyStore }

const keyPrefix = "lnk_"

func NewAuthService(keys APIKeyStore) *AuthService {
	return &AuthService{keys: keys}
}

func generateAPIKey() (plaintext, hash, prefix string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", "", err
	}
	plaintext = keyPrefix + base64.RawURLEncoding.EncodeToString(b)
	hash = hashKey(plaintext)
	prefix = plaintext[:12]
	return plaintext, hash, prefix, nil
}

func hashKey(plaintext string) string {
	h := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(h[:])
}

func (s *AuthService) EnsureAdminKey(ctx context.Context) (plaintext string, created bool, err error) {
	n, err := s.keys.Count(ctx)
	if err != nil {
		return "", false, err
	}
	if n > 0 {
		return "", false, nil
	}
	pt, hash, prefix, err := generateAPIKey()
	if err != nil {
		return "", false, err
	}
	if err := s.keys.Create(ctx, &models.APIKey{KeyHash: hash, Prefix: prefix}); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return "", false, nil
		}
		return "", false, err
	}
	return pt, true, nil
}

func (s *AuthService) EnsureIndexes(ctx context.Context) error {
	return s.keys.EnsureIndexes(ctx)
}

func (s *AuthService) Authenticate(ctx context.Context, plaintext string) (*models.APIKey, error) {
	key, err := s.keys.GetByHash(ctx, hashKey(plaintext))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrInvalidKey
		}
		return nil, err
	}
	_ = s.keys.TouchLastUsed(ctx, key.ID)
	return key, nil
}
