package configs

import (
	"fmt"
)

type Config struct {
	Mongo Mongo
}

type Mongo struct {
	User     string
	Password string
	Host     string
	Port     string
	Timeout  int64
	DB       string
}

func (m *Mongo) GetMongoUrl() string {
	if m.User == "" || m.Password == "" {
		return fmt.Sprintf("mongodb://%s:%s/%s", m.Host, m.Port, m.DB)
	}
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", m.User, m.Password, m.Host, m.Port, m.DB)
}
