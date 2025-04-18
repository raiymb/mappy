package config

import (
	"fmt"
	"time"
)

// This struct mirrors config.yaml / local.yaml (fields are capitalized
// so cleanenv / viper / yaml can populate them).
type Config struct {
	Server    Server    `yaml:"server"`
	Mongo     Mongo     `yaml:"mongo"`
	Redis     Redis     `yaml:"redis"`
	JWT       JWT       `yaml:"jwt"`
	RateLimit RateLimit `yaml:"rateLimit"`
}

type Server struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	Env  string `yaml:"env"  env:"APP_ENV"     env-default:"dev"` // dev | prod
}

func (s Server) ListenAddr() string { return fmt.Sprintf("%s:%d", s.Host, s.Port) }

type Mongo struct {
	URI string `yaml:"uri" env:"MONGO_URI"`
	DB  string `yaml:"db"  env:"MONGO_DB" env-default:"maps_db"`
}

type Redis struct {
	Addr     string `yaml:"addr"     env:"REDIS_ADDR"     env-default:"localhost:6379"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db"       env:"REDIS_DB"       env-default:"0"`
}

type JWT struct {
	Secret     string        `yaml:"secret"     env:"JWT_SECRET"`
	AccessTTL  time.Duration `yaml:"accessTTL"  env:"JWT_ACCESS_TTL"  env-default:"15m"`
	RefreshTTL time.Duration `yaml:"refreshTTL" env:"JWT_REFRESH_TTL" env-default:"720h"`
}

type RateLimit struct {
	Window time.Duration `yaml:"window" env-default:"1m"`
	Max    int           `yaml:"max"    env-default:"120"`
}
