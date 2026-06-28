package setup

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/baydogan/lnk/internal/models"
	"gopkg.in/yaml.v3"
)

const (
	configDirName  = ".lnk"
	serverFileName = "server.yaml"
)

func GetUserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func ConfigDir() (string, error) {
	home, err := GetUserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDirName), nil
}

func ServerConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, serverFileName), nil
}

// ServerConfigExists reports whether server.yaml already exists, along with its path.
func ServerConfigExists() (string, bool, error) {
	path, err := ServerConfigPath()
	if err != nil {
		return "", false, err
	}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return path, false, nil
		}
		return path, false, err
	}
	return path, true, nil
}

func WriteServerConfig(cfg *models.ServerConfig) (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, serverFileName)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", err
	}

	return path, nil
}
