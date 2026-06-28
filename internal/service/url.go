package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/logger"
	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/repository"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{repo: repo}
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

	baseURL := "example.com"

	return &models.ShortenResponse{
		Code:        code,
		ShortURL:    fmt.Sprintf("%s/%s", baseURL, shortCode),
		OriginalURL: url,
	}, nil
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
