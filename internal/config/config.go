package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

var (
	cfg *Config
)

type Config struct {
	Log struct {
		Level string `yaml:"level" envconfig:"LOG_LEVEL"`
	} `yaml:"log"`

	HTTP struct {
		Port               string        `yaml:"port" envconfig:"HTTP_PORT"`
		ReadTimeout        time.Duration `yaml:"read-timeout" envconfig:"HTTP_READ_TIMEOUT"`
		WriteTimeout       time.Duration `yaml:"write-timeout" envconfig:"HTTP_WRITE_TIMEOUT"`
		MaxHeaderMegabytes int           `yaml:"max-header-megabytes" envconfig:"HTTP_MAX_HEADER_MEGABYTES"`
	} `yaml:"http"`

	DB struct {
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
		Name     string `yaml:"name" envconfig:"DB_NAME"`
		User     string `yaml:"user" envconfig:"DB_USER"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
		SSLMode  string `yaml:"ssl-mode" envconfig:"DB_SSL_MODE"`
	} `yaml:"db"`

	Auth struct {
		AccessTokenTTL     time.Duration `yaml:"access-token-ttl" envconfig:"AUTH_ACCESS_TOKEN_TTL"`
		RefreshTokenTTL    time.Duration `yaml:"refresh-token-ttl" envconfig:"AUTH_REFRESH_TOKEN_TTL"`
		RefreshTokenLength int           `yaml:"refresh-token-length" envconfig:"AUTH_REFRESH_TOKEN_LENGTH"`
		PasswordSalt       string        `yaml:"password-salt" envconfig:"AUTH_PASSWORD_SALT"`
		JWT                struct {
			Key string `yaml:"key" envconfig:"AUTH_JWT_KEY"`
		} `yaml:"jwt"`
	} `yaml:"auth"`
}

func LoadConfig(configPath string) *Config {
	if cfg == nil {
		cfg = &Config{}

		cfg.readFile(configPath)
		cfg.readEnv()
	}

	return cfg
}

// File configs with values from configs file
func (c *Config) readFile(path string) {
	f, err := os.Open(path)

	if err != nil {
		processError(err)
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(c)

	if err != nil {
		log.Println(c)
		processError(err)
	}
}

// Read configs with values from env variables
func (c *Config) readEnv() {
	loadFromEnvFile()

	err := envconfig.Process("", c)

	if err != nil {
		processError(err)
	}
}

// Load values from .env file to system
func loadFromEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Debug("Error loading .env file")
	}
}

func processError(err error) {
	log.Error(err)
	os.Exit(2)
}
