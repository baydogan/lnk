package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/baydogan/lnk/domain"
	"gopkg.in/yaml.v3"
)

const (
	configDirName  = ".lnk"
	serverFileName = "server.yaml"
	clientFileName = "config.yaml"

	envConfigPath = "LNK_SERVER_CONFIG"
)

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDirName), nil
}

func ServerConfigPath() (string, error) {
	if p := os.Getenv(envConfigPath); p != "" {
		return p, nil
	}
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, serverFileName), nil
}

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

func ReadServerConfig() (domain.ServerConfig, bool, error) {
	var cfg domain.ServerConfig
	path, err := ServerConfigPath()
	if err != nil {
		return cfg, false, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, false, nil
		}
		return cfg, false, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, false, err
	}
	return cfg, true, nil
}

func WriteServerConfig(cfg *domain.ServerConfig) (string, error) {
	path, err := ServerConfigPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", err
	}
	return path, nil
}

func ClientConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, clientFileName), nil
}

func ReadClientConfig() (domain.ClientConfig, error) {
	var cfg domain.ClientConfig
	path, err := ClientConfigPath()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, domain.ErrNotLoggedIn
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func WriteClientConfig(cfg *domain.ClientConfig) (string, error) {
	path, err := ClientConfigPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", err
	}
	return path, nil
}

func RemoveClientConfig() (string, error) {
	path, err := ClientConfigPath()
	if err != nil {
		return "", err
	}
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return path, domain.ErrNotLoggedIn
		}
		return "", err
	}
	return path, nil
}
