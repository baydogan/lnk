package models

type ServerConfig struct {
	Mode      string `yaml:"mode"`
	MongoURI  string `yaml:"mongo_uri"`
	RedisAddr string `yaml:"redis_addr"`
	Admin     string `yaml:"admin,omitempty"`
}
