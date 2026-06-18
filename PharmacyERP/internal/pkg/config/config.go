package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const defaultConfigPath = "configs/config.local.yaml"

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	Medicare  MedicareConfig  `mapstructure:"medicare"`
	AIService AIServiceConfig `mapstructure:"ai_service"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
}

type ServerConfig struct {
	Mode            string        `mapstructure:"mode"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	SSLMode         string `mapstructure:"ssl_mode"`
	TimeZone        string `mapstructure:"timezone"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime_minutes"`
	LogSQL          bool   `mapstructure:"log_sql"`
}

func (d DatabaseConfig) DSN() string {
	ssl := d.SSLMode
	if ssl == "" {
		ssl = "disable"
	}
	tz := d.TimeZone
	if tz == "" {
		tz = "Asia/Shanghai"
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host,
		d.Port,
		d.User,
		d.Password,
		d.Name,
		ssl,
		tz,
	)
}

type RabbitMQConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	URI       string `mapstructure:"uri"`
	LogQueue  string `mapstructure:"log_queue"`
	Prefetch  int    `mapstructure:"prefetch"`
	Mandatory bool   `mapstructure:"mandatory"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	ExpireTime time.Duration `mapstructure:"expire_time"`
	Issuer     string        `mapstructure:"issuer"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type MedicareConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	BaseURL string        `mapstructure:"base_url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

type AIServiceConfig struct {
	BaseURL string        `mapstructure:"base_url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	if configPath == "" {
		if envPath := os.Getenv("PHARMACY_CONFIG_FILE"); envPath != "" {
			configPath = envPath
		} else {
			configPath = defaultConfigPath
		}
	}

	v.SetConfigFile(configPath)
	v.SetEnvPrefix("PHARMACY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "pharmacy-erp")
	v.SetDefault("app.env", "dev")

	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "10s")
	v.SetDefault("server.shutdown_timeout", "10s")

	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.timezone", "Asia/Shanghai")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.conn_max_lifetime_minutes", 60)
	v.SetDefault("database.log_sql", true)

	v.SetDefault("rabbitmq.enabled", true)
	v.SetDefault("rabbitmq.log_queue", "operation_log_queue")
	v.SetDefault("rabbitmq.prefetch", 20)
	v.SetDefault("rabbitmq.mandatory", false)

	v.SetDefault("jwt.expire_time", "24h")
	v.SetDefault("jwt.issuer", "pharmacy-erp")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "console")

	v.SetDefault("medicare.enabled", false)
	v.SetDefault("medicare.base_url", "http://localhost:8088")
	v.SetDefault("medicare.timeout", "30s")

	v.SetDefault("ai_service.base_url", "http://127.0.0.1:9080")
	v.SetDefault("ai_service.timeout", "120s")
}

func (c *Config) validate() error {
	if c.Database.Host == "" || c.Database.User == "" || c.Database.Name == "" {
		return fmt.Errorf("database config is incomplete")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if c.RabbitMQ.Enabled && c.RabbitMQ.URI == "" {
		return fmt.Errorf("rabbitmq.uri is required when rabbitmq is enabled")
	}
	return nil
}
