package service

import (
	"errors"

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
	return nil, nil
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
