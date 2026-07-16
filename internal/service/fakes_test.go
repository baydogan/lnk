package service

import (
	"context"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type fakeURLStore struct {
	created     []*models.URL
	byKey       map[string]*models.URL
	existing    map[string]bool
	all         []models.URL
	incremented []string
	deleted     []string

	createErr error
	getErr    error
	existsErr error
	listErr   error
	incErr    error
	deleteErr error
}

func newFakeURLStore() *fakeURLStore {
	return &fakeURLStore{byKey: map[string]*models.URL{}, existing: map[string]bool{}}
}

func (f *fakeURLStore) CreateURL(_ context.Context, url *models.URL) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, url)
	return nil
}

func (f *fakeURLStore) GetByCodeOrAlias(_ context.Context, s string) (*models.URL, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if u, ok := f.byKey[s]; ok {
		return u, nil
	}
	return nil, errs.ErrNotFound
}

func (f *fakeURLStore) IncrementClickCount(_ context.Context, code string) error {
	f.incremented = append(f.incremented, code)
	return f.incErr
}

func (f *fakeURLStore) GetURLsByOwner(_ context.Context, ownerID *bson.ObjectID) ([]models.URL, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	if ownerID == nil {
		return f.all, nil
	}
	var out []models.URL
	for _, u := range f.all {
		if u.UserID != nil && *u.UserID == *ownerID {
			out = append(out, u)
		}
	}
	return out, nil
}

func (f *fakeURLStore) CodeExists(_ context.Context, code string) (bool, error) {
	if f.existsErr != nil {
		return false, f.existsErr
	}
	return f.existing[code], nil
}

func (f *fakeURLStore) DeleteByCode(_ context.Context, code string) error {
	f.deleted = append(f.deleted, code)
	return f.deleteErr
}

type fakeKeyStore struct {
	count        int64
	created      []*models.APIKey
	byHash       map[string]*models.APIKey
	touched      []bson.ObjectID
	deletedUsers []bson.ObjectID

	countErr  error
	createErr error
	getErr    error
	touchErr  error
	idxErr    error
	deleteErr error
}

func newFakeKeyStore() *fakeKeyStore {
	return &fakeKeyStore{byHash: map[string]*models.APIKey{}}
}

func (f *fakeKeyStore) Count(_ context.Context) (int64, error) {
	return f.count, f.countErr
}

func (f *fakeKeyStore) Create(_ context.Context, k *models.APIKey) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = append(f.created, k)
	return nil
}

func (f *fakeKeyStore) EnsureIndexes(_ context.Context) error {
	return f.idxErr
}

func (f *fakeKeyStore) GetByHash(_ context.Context, hash string) (*models.APIKey, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if k, ok := f.byHash[hash]; ok {
		return k, nil
	}
	return nil, errs.ErrNotFound
}

func (f *fakeKeyStore) TouchLastUsed(_ context.Context, id bson.ObjectID) error {
	f.touched = append(f.touched, id)
	return f.touchErr
}

func (f *fakeKeyStore) DeleteByUserID(_ context.Context, id bson.ObjectID) error {
	f.deletedUsers = append(f.deletedUsers, id)
	return f.deleteErr
}

type fakeUserStore struct {
	created []*models.User
	byName  map[string]*models.User
	deleted []string

	createErr error
	getErr    error
	listErr   error
}

func newFakeUserStore() *fakeUserStore {
	return &fakeUserStore{byName: map[string]*models.User{}}
}

func (f *fakeUserStore) Create(_ context.Context, u *models.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	u.ID = bson.NewObjectID()
	f.created = append(f.created, u)
	f.byName[u.Username] = u
	return nil
}

func (f *fakeUserStore) GetByUsername(_ context.Context, username string) (*models.User, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if u, ok := f.byName[username]; ok {
		return u, nil
	}
	return nil, errs.ErrNotFound
}

func (f *fakeUserStore) GetByID(_ context.Context, id bson.ObjectID) (*models.User, error) {
	for _, u := range f.byName {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errs.ErrNotFound
}

func (f *fakeUserStore) List(_ context.Context) ([]models.User, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	out := make([]models.User, 0, len(f.created))
	for _, u := range f.created {
		out = append(out, *u)
	}
	return out, nil
}

func (f *fakeUserStore) EnsureIndexes(_ context.Context) error { return nil }

func (f *fakeUserStore) DeleteByUsername(_ context.Context, username string) error {
	if _, ok := f.byName[username]; !ok {
		return errs.ErrNotFound
	}
	delete(f.byName, username)
	f.deleted = append(f.deleted, username)
	return nil
}
