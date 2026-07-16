//go:build integration

package mongo

import (
	"context"
	"log"
	"os"
	"testing"

	driver "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

var testDB *driver.Database

func TestMain(m *testing.M) {
	ctx := context.Background()

	ctr, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		log.Fatalf("start mongo container: %v", err)
	}

	uri, err := ctr.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("connection string: %v", err)
	}
	db, err := Connect(ctx, uri)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	testDB = db

	code := m.Run()

	_ = ctr.Terminate(ctx)
	os.Exit(code)
}

func clearCollection(t *testing.T, name string) {
	t.Helper()
	if err := testDB.Collection(name).Drop(context.Background()); err != nil {
		t.Fatalf("drop %s: %v", name, err)
	}
}
