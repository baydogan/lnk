package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/logger"
	"github.com/baydogan/lnk/internal/models"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type URLService struct {
	repo    URLStore
	baseURL string
}

func NewURLService(repo URLStore, baseURL string) *URLService {
	return &URLService{repo: repo, baseURL: baseURL}
}

func (s *URLService) ShortenURL(ctx context.Context, req *models.ShortenRequest, ownerID *bson.ObjectID) (*models.ShortenResponse, error) {
	url := strings.TrimSpace(req.URL)
	if url == "" {
		return nil, errs.ErrInvalidURL
	}

	if !strings.Contains(url, "//") {
		url = "https://" + url
	}

	if req.Alias != "" {
		exists, err := s.repo.CodeExists(ctx, req.Alias)

		if err != nil {
			return nil, err
		}

		if exists {
			return nil, errs.ErrAliasExists
		}
	}

	expiresAt, err := parseExpiry(req.Expires)
	if err != nil {
		return nil, err
	}

	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return nil, err
	}

	var alias *string
	if req.Alias != "" {
		alias = &req.Alias
	}

	urlModel := &models.URL{
		Code:        code,
		OriginalURL: url,
		Alias:       alias,
		ExpiresAt:   expiresAt,
		UserID:      ownerID,
	}

	if err := s.repo.CreateURL(ctx, urlModel); err != nil {
		return nil, err
	}

	shortCode := code
	if alias != nil {
		shortCode = *alias
	}

	logger.Info().
		Str("code", code).
		Str("original_url", url).
		Msg("URL shortened")

	return &models.ShortenResponse{
		Code:        code,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), shortCode),
		OriginalURL: url,
	}, nil
}

func (s *URLService) ResolveURL(ctx context.Context, codeOrAlias string) (string, error) {
	u, err := s.repo.GetByCodeOrAlias(ctx, codeOrAlias)
	if err != nil {
		return "", err
	}
	if u.ExpiresAt != nil && time.Now().After(*u.ExpiresAt) {
		return "", errs.ErrNotFound
	}
	if err := s.repo.IncrementClickCount(ctx, u.Code); err != nil {
		logger.Error().Err(err).Str("code", u.Code).Msg("click increment failed")
	}
	return u.OriginalURL, nil
}

func (s *URLService) ListURLs(ctx context.Context, ownerID *bson.ObjectID) ([]models.URLResponse, error) {
	urls, err := s.repo.GetURLsByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]models.URLResponse, 0, len(urls))
	for i := range urls {
		r := urls[i].ToResponse()
		r.ShortURL = s.shortURL(&urls[i])
		out = append(out, r)
	}
	return out, nil
}

func (s *URLService) GetURL(ctx context.Context, codeOrAlias string, ownerID *bson.ObjectID) (*models.URLResponse, error) {
	u, err := s.repo.GetByCodeOrAlias(ctx, codeOrAlias)
	if err != nil {
		return nil, err
	}
	if !owns(u, ownerID) {
		return nil, errs.ErrNotFound
	}
	r := u.ToResponse()
	r.ShortURL = s.shortURL(u)
	return &r, nil
}

func owns(u *models.URL, ownerID *bson.ObjectID) bool {
	if ownerID == nil {
		return true
	}
	return u.UserID != nil && *u.UserID == *ownerID
}

func (s *URLService) shortURL(u *models.URL) string {
	code := u.Code
	if u.Alias != nil {
		code = *u.Alias
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), code)
}

func (s *URLService) DeleteURL(ctx context.Context, codeOrAlias string, ownerID *bson.ObjectID) error {
	if ownerID != nil {
		u, err := s.repo.GetByCodeOrAlias(ctx, codeOrAlias)
		if err != nil {
			return err
		}
		if !owns(u, ownerID) {
			return errs.ErrNotFound
		}
	}
	return s.repo.DeleteByCode(ctx, codeOrAlias)
}

func (s *URLService) generateUniqueCode(ctx context.Context) (string, error) {
	for range 5 {
		code, err := gonanoid.New(7)
		if err != nil {
			return "", err
		}

		exists, err := s.repo.CodeExists(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique code after 5 attempts")
}

func parseExpiry(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	var d time.Duration
	switch last := s[len(s)-1]; last {
	case 'd', 'w':
		n, err := strconv.Atoi(s[:len(s)-1])
		if err != nil || n <= 0 {
			return nil, errs.ErrExpireFormat
		}
		unit := 24 * time.Hour
		if last == 'w' {
			unit = 7 * 24 * time.Hour
		}
		d = time.Duration(n) * unit
	default:
		parsed, err := time.ParseDuration(s)
		if err != nil || parsed <= 0 {
			return nil, errs.ErrExpireFormat
		}
		d = parsed
	}

	t := time.Now().Add(d)
	return &t, nil
}
