package config

import (
	"fmt"
	"os"
)

type Db struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
	HttpUrl  string `json:"http_url"`
}

type Config struct {
	Dgraph Db
	Port   string
	Auth   string
}

func LoadConfig() *Config {
	db := Db{
		Url:      os.Getenv("DB_URL"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		HttpUrl:  os.Getenv("DB_HTTP_URL"),
	}
	if db.Url == "" {
		db.Url = "localhost:9080"
	}
	if db.HttpUrl == "" {
		db.HttpUrl = "localhost:8080"
	}
	if db.User == "" {
		db.User = ""
	}
	if db.Password == "" {
		db.Password = ""
	}
	auth := os.Getenv("AUTH_API_SECRET")
	return &Config{
		Dgraph: db,
		Port:   "9000",
		Auth:   auth,
	}
}

func (c *Config) Address() string {
	return fmt.Sprintf(":%s", c.Port)
}
