package user

import (
	"context"
	"errors"
	"strings"

	"github.com/baydogan/lnk/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Repository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id bson.ObjectID) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
	EnsureIndexes(ctx context.Context) error
	DeleteByUsername(ctx context.Context, username string) error
}

type KeyRepository interface {
	Create(ctx context.Context, k *domain.APIKey) error
	DeleteByUserID(ctx context.Context, userID bson.ObjectID) error
}

type Service struct {
	users Repository
	keys  KeyRepository
}

func NewService(users Repository, keys KeyRepository) *Service {
	return &Service{users: users, keys: keys}
}

func (s *Service) CreateUser(ctx context.Context, username, role string) (*domain.User, string, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, "", domain.ErrInvalidUsername
	}
	if role != domain.RoleAdmin && role != domain.RoleUser {
		return nil, "", domain.ErrInvalidRole
	}

	user := &domain.User{Username: username, Role: role}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, "", err
	}

	plaintext, err := s.issueKey(ctx, &user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, plaintext, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.users.List(ctx)
}

func (s *Service) GetUser(ctx context.Context, id bson.ObjectID) (*domain.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *Service) DeleteUser(ctx context.Context, username string) error {
	username = strings.TrimSpace(username)
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.Role == domain.RoleAdmin {
		return domain.ErrCannotDeleteAdmin
	}
	if err := s.keys.DeleteByUserID(ctx, user.ID); err != nil {
		return err
	}
	return s.users.DeleteByUsername(ctx, username)
}

func (s *Service) EnsureIndexes(ctx context.Context) error {
	return s.users.EnsureIndexes(ctx)
}

func (s *Service) EnsureAdmin(ctx context.Context, username string) (plaintext string, created bool, err error) {
	if _, err := s.users.GetByUsername(ctx, username); err == nil {
		return "", false, nil
	} else if !errors.Is(err, domain.ErrNotFound) {
		return "", false, err
	}

	_, pt, err := s.CreateUser(ctx, username, domain.RoleAdmin)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			return "", false, nil
		}
		return "", false, err
	}
	return pt, true, nil
}

func (s *Service) issueKey(ctx context.Context, userID *bson.ObjectID) (string, error) {
	plaintext, hash, prefix, err := domain.GenerateAPIKey()
	if err != nil {
		return "", err
	}
	if err := s.keys.Create(ctx, &domain.APIKey{KeyHash: hash, Prefix: prefix, UserID: userID}); err != nil {
		return "", err
	}
	return plaintext, nil
}
