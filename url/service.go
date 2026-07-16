package url

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/baydogan/lnk/domain"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repository interface {
	CreateURL(ctx context.Context, url *domain.URL) error
	GetByCodeOrAlias(ctx context.Context, s string) (*domain.URL, error)
	IncrementClickCount(ctx context.Context, code string) error
	GetURLsByOwner(ctx context.Context, ownerID *bson.ObjectID) ([]domain.URL, error)
	CodeExists(ctx context.Context, code string) (bool, error)
	DeleteByCode(ctx context.Context, code string) error
}

type Service struct {
	repo    Repository
	baseURL string
}

func NewService(repo Repository, baseURL string) *Service {
	return &Service{repo: repo, baseURL: baseURL}
}

func (s *Service) ShortenURL(ctx context.Context, req *domain.ShortenRequest, ownerID *bson.ObjectID) (*domain.ShortenResponse, error) {
	url := strings.TrimSpace(req.URL)
	if url == "" {
		return nil, domain.ErrInvalidURL
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
			return nil, domain.ErrAliasExists
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

	urlModel := &domain.URL{
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

	return &domain.ShortenResponse{
		Code:        code,
		ShortURL:    fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), shortCode),
		OriginalURL: url,
	}, nil
}

func (s *Service) ResolveURL(ctx context.Context, codeOrAlias string) (string, error) {
	u, err := s.repo.GetByCodeOrAlias(ctx, codeOrAlias)
	if err != nil {
		return "", err
	}
	if u.ExpiresAt != nil && time.Now().After(*u.ExpiresAt) {
		return "", domain.ErrNotFound
	}
	_ = s.repo.IncrementClickCount(ctx, u.Code)
	return u.OriginalURL, nil
}

func (s *Service) ListURLs(ctx context.Context, ownerID *bson.ObjectID) ([]domain.URLResponse, error) {
	urls, err := s.repo.GetURLsByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.URLResponse, 0, len(urls))
	for i := range urls {
		r := urls[i].ToResponse()
		r.ShortURL = s.shortURL(&urls[i])
		out = append(out, r)
	}
	return out, nil
}

func (s *Service) GetURL(ctx context.Context, codeOrAlias string, ownerID *bson.ObjectID) (*domain.URLResponse, error) {
	u, err := s.repo.GetByCodeOrAlias(ctx, codeOrAlias)
	if err != nil {
		return nil, err
	}
	if !owns(u, ownerID) {
		return nil, domain.ErrNotFound
	}
	r := u.ToResponse()
	r.ShortURL = s.shortURL(u)
	return &r, nil
}

func owns(u *domain.URL, ownerID *bson.ObjectID) bool {
	if ownerID == nil {
		return true
	}
	return u.UserID != nil && *u.UserID == *ownerID
}

func (s *Service) shortURL(u *domain.URL) string {
	code := u.Code
	if u.Alias != nil {
		code = *u.Alias
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), code)
}

func (s *Service) DeleteURL(ctx context.Context, codeOrAlias string, ownerID *bson.ObjectID) error {
	if ownerID != nil {
		u, err := s.repo.GetByCodeOrAlias(ctx, codeOrAlias)
		if err != nil {
			return err
		}
		if !owns(u, ownerID) {
			return domain.ErrNotFound
		}
	}
	return s.repo.DeleteByCode(ctx, codeOrAlias)
}

func (s *Service) generateUniqueCode(ctx context.Context) (string, error) {
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
			return nil, domain.ErrExpireFormat
		}
		unit := 24 * time.Hour
		if last == 'w' {
			unit = 7 * 24 * time.Hour
		}
		d = time.Duration(n) * unit
	default:
		parsed, err := time.ParseDuration(s)
		if err != nil || parsed <= 0 {
			return nil, domain.ErrExpireFormat
		}
		d = parsed
	}

	t := time.Now().Add(d)
	return &t, nil
}
