package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/logger"
	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/repository"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type URLService struct {
	repo    *repository.URLRepository
	baseURL string
}

func NewURLService(repo *repository.URLRepository, baseURL string) *URLService {
	return &URLService{repo: repo, baseURL: baseURL}
}

func (s *URLService) ShortenURL(req *models.ShortenRequest) (*models.ShortenResponse, error) {
	url := strings.TrimSpace(req.URL)
	if url == "" {
		return nil, errs.ErrInvalidURL
	}

	if !strings.Contains(url, "//") {
		url = "https://" + url
	}

	if req.Alias != "" {
		exists, err := s.repo.CodeExists(req.Alias)

		if err != nil {
			return nil, err
		}

		if exists {
			return nil, errs.ErrAliasExists
		}
	}

	code, err := s.generateUniqueCode()

	if err != nil {
		return nil, err
	}

	//TODO parse
	//var expiresAt time.Time

	var alias *string
	if req.Alias != "" {
		alias = &req.Alias
	}

	urlModel := &models.URL{
		Code:        code,
		OriginalURL: url,
		Alias:       alias,
	}

	if err := s.repo.CreateURL(urlModel); err != nil {
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

func (s *URLService) ResolveURL(codeOrAlias string) (string, error) {
	u, err := s.repo.GetByCodeOrAlias(codeOrAlias)
	if err != nil {
		return "", err
	}
	if u.ExpiresAt != nil && time.Now().After(*u.ExpiresAt) {
		return "", errs.ErrNotFound
	}
	if err := s.repo.IncrementClickCount(u.Code); err != nil {
		logger.Error().Err(err).Str("code", u.Code).Msg("click increment failed")
	}
	return u.OriginalURL, nil
}

func (s *URLService) ListURLs() ([]models.URLResponse, error) {
	urls, err := s.repo.GetAllURLs()
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

func (s *URLService) GetURL(codeOrAlias string) (*models.URLResponse, error) {
	u, err := s.repo.GetByCodeOrAlias(codeOrAlias)
	if err != nil {
		return nil, err
	}
	r := u.ToResponse()
	r.ShortURL = s.shortURL(u)
	return &r, nil
}

func (s *URLService) shortURL(u *models.URL) string {
	code := u.Code
	if u.Alias != nil {
		code = *u.Alias
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), code)
}

func (s *URLService) DeleteURL(codeOrAlias string) error {
	return s.repo.DeleteByCode(codeOrAlias)
}

func (s *URLService) generateUniqueCode() (string, error) {
	for range 5 {
		code, err := gonanoid.New(7)
		if err != nil {
			return "", err
		}

		exists, err := s.repo.CodeExists(code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique code after 5 attempts")
}
