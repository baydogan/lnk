package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/baydogan/lnk/internal/errs"
	"github.com/baydogan/lnk/internal/models"
)

func TestServerConfigRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.yaml")
	t.Setenv(envConfigPath, path)

	want := &models.ServerConfig{
		Mode:      "single",
		BaseURL:   "http://localhost:8080",
		MongoURI:  "mongodb://localhost:27017",
		RedisAddr: "localhost:6379",
		Admin:     "me",
	}
	if _, err := WriteServerConfig(want); err != nil {
		t.Fatalf("WriteServerConfig: %v", err)
	}

	got, ok, err := ReadServerConfig()
	if err != nil {
		t.Fatalf("ReadServerConfig: %v", err)
	}
	if !ok {
		t.Fatal("ReadServerConfig ok = false, want true")
	}
	if got != *want {
		t.Fatalf("roundtrip mismatch: got %+v want %+v", got, *want)
	}
}

func TestReadServerConfigMissing(t *testing.T) {
	t.Setenv(envConfigPath, filepath.Join(t.TempDir(), "absent.yaml"))
	cfg, ok, err := ReadServerConfig()
	if err != nil {
		t.Fatalf("ReadServerConfig missing: unexpected error %v", err)
	}
	if ok {
		t.Fatal("ReadServerConfig ok = true for missing file")
	}
	if cfg != (models.ServerConfig{}) {
		t.Fatalf("ReadServerConfig missing returned %+v, want zero", cfg)
	}
}

func TestClientConfigRoundtripAndPerms(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	want := &models.ClientConfig{Server: "http://localhost:8080", APIKey: "lnk_secret"}
	path, err := WriteClientConfig(want)
	if err != nil {
		t.Fatalf("WriteClientConfig: %v", err)
	}

	got, err := ReadClientConfig()
	if err != nil {
		t.Fatalf("ReadClientConfig: %v", err)
	}
	if got != *want {
		t.Fatalf("roundtrip mismatch: got %+v want %+v", got, *want)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat file: %v", err)
	}
	if perm := fi.Mode().Perm(); perm != 0o600 {
		t.Fatalf("file perm = %o, want 600", perm)
	}
	di, err := os.Stat(filepath.Dir(path))
	if err != nil {
		t.Fatalf("stat dir: %v", err)
	}
	if perm := di.Mode().Perm(); perm != 0o700 {
		t.Fatalf("dir perm = %o, want 700", perm)
	}
}

func TestReadClientConfigNotLoggedIn(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if _, err := ReadClientConfig(); !errors.Is(err, errs.ErrNotLoggedIn) {
		t.Fatalf("ReadClientConfig err = %v, want ErrNotLoggedIn", err)
	}
}

func TestRemoveClientConfig(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	if _, err := WriteClientConfig(&models.ClientConfig{Server: "s", APIKey: "k"}); err != nil {
		t.Fatalf("WriteClientConfig: %v", err)
	}
	if _, err := RemoveClientConfig(); err != nil {
		t.Fatalf("RemoveClientConfig first call: %v", err)
	}
	if _, err := RemoveClientConfig(); !errors.Is(err, errs.ErrNotLoggedIn) {
		t.Fatalf("RemoveClientConfig second call err = %v, want ErrNotLoggedIn", err)
	}
}
