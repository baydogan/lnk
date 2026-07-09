//go:build integration

package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/baydogan/lnk/internal/database"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

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
	if err := database.Connect(uri); err != nil {
		log.Fatalf("connect: %v", err)
	}

	code := m.Run()

	_ = ctr.Terminate(ctx)
	os.Exit(code)
}

func clearCollection(t *testing.T, name string) {
	t.Helper()
	if err := database.Collection(name).Drop(context.Background()); err != nil {
		t.Fatalf("drop %s: %v", name, err)
	}
}
