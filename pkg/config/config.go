package config

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	publicCertsTTL = 24
)

type singleton struct{}

var once sync.Once

var appConfig Config

type loader interface {
	Load(viper.Viper) (*viper.Viper, error)
}

// Config contain configuration of db for migrator
// config var < env < command flag
type Config struct {
	ServiceName    string
	Version        string
	BaseURL        string
	Port           string
	Env            string
	AllowedOrigins string
	AccessTokenTTL time.Duration
	JWTSecret      []byte
	DBHost         string
	DBPort         string
	DBUser         string
	DBName         string
	DBPass         string
	DBSSLMode      string
	SentryDSN      string
}

// generateConfigFromViper generate config from viper data
func generateConfigFromViper(v *viper.Viper) Config {
	tokenTTLInDay := v.GetInt("ACCESS_TOKEN_TTL")
	if tokenTTLInDay == 0 {
		tokenTTLInDay = 7
	}

	return Config{
		ServiceName: v.GetString("SERVICE_NAME"),
		Port:        v.GetString("PORT"),
		BaseURL:     v.GetString("BASE_URL"),
		Version:     v.GetString("VERSION"),
		Env:         v.GetString("ENV"),

		AllowedOrigins: v.GetString("ALLOWED_ORIGINS"),

		DBHost:    v.GetString("DB_HOST"),
		DBPort:    v.GetString("DB_PORT"),
		DBUser:    v.GetString("DB_USER"),
		DBName:    v.GetString("DB_NAME"),
		DBPass:    v.GetString("DB_PASS"),
		DBSSLMode: v.GetString("DB_SSL_MODE"),

		AccessTokenTTL: time.Hour * 24 * time.Duration(tokenTTLInDay),
		JWTSecret:      []byte(v.GetString("JWT_SECRET")),

		SentryDSN: v.GetString("SENTRY_DSN"),
	}
}

// GetDBURL in config
func (c *Config) GetDBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

// GetCORS in config
func (c *Config) GetCORS() []string {
	cors := strings.Split(c.AllowedOrigins, ";")
	rs := []string{}
	for idx := range cors {
		itm := cors[idx]
		if strings.TrimSpace(itm) != "" {
			rs = append(rs, itm)
		}
	}
	return rs
}

// GetShutdownTimeout get shutdown time out
func (c *Config) GetShutdownTimeout() time.Duration {
	return 10 * time.Second
}

var envDir = "."
var envFileName = ".env"

func GetConfig() Config {
	once.Do(func() {
		loaders := []loader{}

		fileLoader := NewFileLoader(envFileName, envDir)
		loaders = append(loaders, fileLoader)

		loaders = append(loaders, NewENVLoader())

		v := viper.New()
		v.SetDefault("PORT", "8100")
		v.SetDefault("ENV", "local")
		v.SetDefault("DB_HOST", "127.0.0.1")
		v.SetDefault("DB_PORT", "54321")
		v.SetDefault("DB_USER", "postgres")
		v.SetDefault("DB_PASS", "postgres")
		v.SetDefault("DB_NAME", "db_local")
		v.SetDefault("VERSION", "0.0.0")
		v.SetDefault("JWT_SECRET", "sample_secret")

		for idx := range loaders {
			newV, err := loaders[idx].Load(*v)

			if err == nil {
				v = newV
			}
		}

		appConfig = generateConfigFromViper(v)
	})

	return appConfig
}
