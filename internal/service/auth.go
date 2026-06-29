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
	"github.com/baydogan/lnk/internal/repository"
)

type AuthService struct{ keys *repository.APIKeyRepository }

const keyPrefix = "lnk_"

func NewAuthService(keys *repository.APIKeyRepository) *AuthService {
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

func (s *AuthService) EnsureAdminKey() (plaintext string, created bool, err error) {
	n, err := s.keys.Count()
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
	if err := s.keys.Create(&models.APIKey{KeyHash: hash, Prefix: prefix}); err != nil {
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
