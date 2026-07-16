package service

import (
	"context"
	"errors"
	"strings"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService struct {
	users UserStore
	keys  APIKeyStore
}

func NewUserService(users UserStore, keys APIKeyStore) *UserService {
	return &UserService{users: users, keys: keys}
}

func (s *UserService) CreateUser(ctx context.Context, username, role string) (*models.User, string, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, "", errs.ErrInvalidUsername
	}
	if role != models.RoleAdmin && role != models.RoleUser {
		return nil, "", errs.ErrInvalidRole
	}

	user := &models.User{Username: username, Role: role}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, "", err
	}

	plaintext, err := s.issueKey(ctx, &user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, plaintext, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.users.List(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, username string) error {
	username = strings.TrimSpace(username)
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user.Role == models.RoleAdmin {
		return errs.ErrCannotDeleteAdmin
	}
	if err := s.keys.DeleteByUserID(ctx, user.ID); err != nil {
		return err
	}
	return s.users.DeleteByUsername(ctx, username)
}

func (s *UserService) EnsureIndexes(ctx context.Context) error {
	return s.users.EnsureIndexes(ctx)
}

func (s *UserService) EnsureAdmin(ctx context.Context, username string) (plaintext string, created bool, err error) {
	if _, err := s.users.GetByUsername(ctx, username); err == nil {
		return "", false, nil
	} else if !errors.Is(err, errs.ErrNotFound) {
		return "", false, err
	}

	_, pt, err := s.CreateUser(ctx, username, models.RoleAdmin)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			return "", false, nil
		}
		return "", false, err
	}
	return pt, true, nil
}

func (s *UserService) issueKey(ctx context.Context, userID *bson.ObjectID) (string, error) {
	plaintext, hash, prefix, err := generateAPIKey()
	if err != nil {
		return "", err
	}
	if err := s.keys.Create(ctx, &models.APIKey{KeyHash: hash, Prefix: prefix, UserID: userID}); err != nil {
		return "", err
	}
	return plaintext, nil
}
