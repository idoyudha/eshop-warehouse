package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"log"`
		PostgreSQL
		AuthService
		Redis
		Kafka
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	PostgreSQL struct {
		URL         string `env-required:"true" env:"POSTGRESQL_URL"`
		ConnTimeout int    `env-required:"true" env:"POSTGRESQL_CONN_TIMEOUT"`
		ConnAttemps int    `env-required:"true" env:"POSTGRESQL_CONN_ATTEMPS"`
		MaxPoolSize int    `env-required:"true" env:"POSTGRESQL_MAX_POOL_SIZE"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	AuthService struct {
		BaseURL string `env-required:"true" env:"AUTH_SERVICE"`
	}

	Redis struct {
		RedisURL      string `env-required:"true" env:"REDIS_URL"`
		RedisPassword string `env-required:"true" env:"REDIS_PASSWORD"`
	}

	Kafka struct {
		Broker string `env-required:"true" env:"KAFKA_BROKER"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
