package conf

import "os"

type Db struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Config struct {
	Dgraph Db
}

func LoadConfig() Config {
	db := Db{
		Url:      os.Getenv("DB_URL"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
	}
	if db.Url == "" {
		db.Url = "localhost:9080"
	}
	if db.User == "" {
		db.User = ""
	}
	if db.Password == "" {
		db.Password = ""
	}
	return Config{
		Dgraph: db,
	}
}
