package models

const (
	ModeSingle = "single"
	ModeMulti  = "multi"
)

type ServerConfig struct {
	Mode      string `yaml:"mode"`
	BaseURL   string `yaml:"base_url"`
	MongoURI  string `yaml:"mongo_uri"`
	RedisAddr string `yaml:"redis_addr"`
	Admin     string `yaml:"admin,omitempty"`
}
