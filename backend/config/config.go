package config

import (
	"backend/pkg/logger"
	"embed"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AppName         string              `mapstructure:"app_name" yaml:"app_name"`
	LogLevel        string              `mapstructure:"log_level" yaml:"log_level"`
	DB              DBConfig            `mapstructure:"database" yaml:"database"`
	RabbitMQ        RabbitMQConfig      `mapstructure:"rabbitmq" yaml:"rabbitmq"`
	WebServer       ServerConfig        `mapstructure:"web_server" yaml:"web_server"`
	Frontend        FrontendConfig      `mapstructure:"frontend" yaml:"frontend"`
	Register        RegisterConfig      `mapstructure:"register" yaml:"register"`
	ResetPassword   ResetPasswordConfig `mapstructure:"reset_password" yaml:"reset_password"`
	EmailChange     EmailChangeConfig   `mapstructure:"email_change" yaml:"email_change"`
	Email           EmailConfig         `mapstructure:"email" yaml:"email"`
	DefaultLanguage string              `mapstructure:"default_language" yaml:"default_language"`
}

type EmailChangeConfig struct {
	ExpirationDays int `mapstructure:"expiration_days" yaml:"expiration_days"`
}

type RegisterConfig struct {
	Enabled              bool   `mapstructure:"enabled" yaml:"enabled"`
	ConfirmationEndpoint string `mapstructure:"confirmation_endpoint" yaml:"confirmation_endpoint"`
	ExpirationDays       int    `mapstructure:"expiration_days" yaml:"expiration_days"`
}

type ResetPasswordConfig struct {
	Enabled        bool `mapstructure:"enabled" yaml:"enabled"`
	ExpirationDays int  `mapstructure:"expiration_days" yaml:"expiration_days"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	URL      string `mapstructure:"url" yaml:"url"`
}

type DBConfig struct {
	DSN    string `mapstructure:"dsn" yaml:"dsn"`
	User   string `mapstructure:"user" yaml:"user"`
	Pass   string `mapstructure:"password" yaml:"password"`
	Port   int    `mapstructure:"port" yaml:"port"`
	Host   string `mapstructure:"host" yaml:"host"`
	DbName string `mapstructure:"dbname" yaml:"dbname"`
}

type ServerConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	HTTPPort int    `mapstructure:"http_port" yaml:"http_port"`
}
type FrontendConfig struct {
	BaseURL              string `mapstructure:"base_url" yaml:"base_url"`
	ConfirmationEndpoint string `mapstructure:"confirmation_endpoint" yaml:"confirmation_endpoint"`
}
type FeaturesConfig struct {
	Register      bool `json:"register" yaml:"register"`
	ResetPassword bool `json:"reset_password" yaml:"reset_password"`
}
type EmailConfig struct {
	SMTPHost string `mapstructure:"smtp_host" yaml:"smtp_host"`
	SMTPPort string `mapstructure:"smtp_port" yaml:"smtp_port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	From     string `mapstructure:"from" yaml:"from"`
}

var configInstance *Config

//go:embed *.yaml
var configFiles embed.FS

func GetConfig() (*Config, error) {
	if configInstance == nil {
		cfg, err := loadConfig()
		if err != nil {
			return nil, err
		}
		configInstance = cfg
	}
	return configInstance, nil
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	v := viper.New()
	v.SetConfigType("yaml")
	setDefaults(v)

	files := []string{"config.yaml"}
	if env := os.Getenv("APP_ENV"); env != "" {
		files = append(files, fmt.Sprintf("config_%s.yaml", env))
	}
	for _, f := range files {
		data, err := configFiles.ReadFile(f)
		if err != nil {
			fmt.Printf("File %s does not exist: %v\n", f, err)
			continue
		}
		if err := v.MergeConfigMap(readYamlToMap(data)); err != nil {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	setConfigByEnv(&cfg)
	if cfg.DB.DSN == "" {
		cfg.DB.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.DbName)
	}
	if cfg.RabbitMQ.URL == "" {
		// "amqp://guest:guest@localhost:5672/"
		cfg.RabbitMQ.URL = fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
	}
	return &cfg, nil
}
func readYamlToMap(data []byte) map[string]interface{} {
	m := make(map[string]interface{})
	_ = yaml.Unmarshal(data, &m)
	return m
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app_name", "WebApp")
	v.SetDefault("log_level", "info")
	v.SetDefault("web_server.http_port", 8080)
	v.SetDefault("web_server.host", "")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.host", "headless-db")
	v.SetDefault("rabbitmq.user", "guest")
	v.SetDefault("rabbitmq.password", "guest")
	v.SetDefault("rabbitmq.port", 5672)
	v.SetDefault("rabbitmq.host", "rabbitmq")
}
func setConfigByEnv(cfg *Config) {
	setGeneralConfigByEnv(cfg)
	setWebServerConfigByEnv(cfg)
	setDatabaseConfigByEnv(cfg)
	setRabbitMQConfigByEnv(cfg)
	setFrontendConfigByEnv(cfg)
	setResetPasswordConfigByEnv(cfg)
	setRegisterConfigByEnv(cfg)
	setEmailConfigByEnv(cfg)
	setEmailChangeConfigByEnv(cfg)
}

func setEmailConfigByEnv(cfg *Config) {
	if smtpHost := os.Getenv("SMTP_HOST"); smtpHost != "" {
		cfg.Email.SMTPHost = smtpHost
	}
	if smtpPort := os.Getenv("SMTP_PORT"); smtpPort != "" {
		cfg.Email.SMTPPort = smtpPort
	}
	if username := os.Getenv("SMTP_USERNAME"); username != "" {
		cfg.Email.Username = username
	}
	if password := os.Getenv("SMTP_PASSWORD"); password != "" {
		cfg.Email.Password = password
	}
	if from := os.Getenv("SMTP_FROM"); from != "" {
		cfg.Email.From = from
	}
}

func setRegisterConfigByEnv(cfg *Config) {
	if enabled := os.Getenv("REGISTER_ENABLED"); enabled != "" {
		cfg.Register.Enabled = enabled == "true"
	}
	if confirmationEndpoint := os.Getenv("REGISTER_CONFIRMATION_ENDPOINT"); confirmationEndpoint != "" {
		cfg.Register.ConfirmationEndpoint = confirmationEndpoint
	}
	if expirationDays := os.Getenv("REGISTER_EXPIRATION_DAYS"); expirationDays != "" {
		intDays, err := strconv.Atoi(expirationDays)
		if err != nil {
			logger.Error("Invalid REGISTER_EXPIRATION_DAYS value: %v; setting to default", err)
			intDays = 1 // default expiration days
		}
		cfg.Register.ExpirationDays = intDays
	}
}
func setResetPasswordConfigByEnv(cfg *Config) {
	if enabled := os.Getenv("RESET_PASSWORD_ENABLED"); enabled != "" {
		cfg.ResetPassword.Enabled = enabled == "true"
	}
	if expirationDays := os.Getenv("RESET_PASSWORD_EXPIRATION_DAYS"); expirationDays != "" {
		intDays, err := strconv.Atoi(expirationDays)
		if err != nil {
			logger.Error("Invalid RESET_PASSWORD_EXPIRATION_DAYS value: %v; setting to default", err)
			intDays = 1 // default expiration days
		}
		cfg.ResetPassword.ExpirationDays = intDays
	}
}
func setEmailChangeConfigByEnv(cfg *Config) {
	if expirationDays := os.Getenv("EMAIL_CHANGE_EXPIRATION_DAYS"); expirationDays != "" {
		intDays, err := strconv.Atoi(expirationDays)
		if err != nil {

			logger.Error("Invalid EMAIL_CHANGE_EXPIRATION_DAYS value: %v; setting to default", err)
			intDays = 1 // default expiration days
		}
		cfg.EmailChange.ExpirationDays = intDays
	}
}

func setWebServerConfigByEnv(cfg *Config) {
	if host := os.Getenv("BACKEND_HOST"); host != "" {
		cfg.WebServer.Host = host
	}
	if httpPort := os.Getenv("WEBSERVER_PORT"); httpPort != "" {
		intPort, err := strconv.Atoi(httpPort)
		if err != nil {
			logger.Error("Invalid DB_PORT value: %v; setting to default", err)
			intPort = 8080 // default HTTP port
		}
		cfg.WebServer.HTTPPort = intPort
	}
}

func setDatabaseConfigByEnv(cfg *Config) {
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.DB.User = dbUser
	}
	if dbPass := os.Getenv("DB_PASS"); dbPass != "" {
		cfg.DB.Pass = dbPass
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.DB.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		intPort, err := strconv.Atoi(dbPort)
		if err != nil {
			logger.Error("Invalid DB_PORT value: %v; setting to default", err)
			intPort = 3306 // default MySQL port
		}
		cfg.DB.Port = intPort
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.DB.DbName = dbName
	}
}

func setRabbitMQConfigByEnv(cfg *Config) {
	if rabbitMQHost := os.Getenv("RABBITMQ_HOST"); rabbitMQHost != "" {
		cfg.RabbitMQ.Host = rabbitMQHost
	}
	if rabbitMQPort := os.Getenv("RABBITMQ_PORT"); rabbitMQPort != "" {
		intPort, err := strconv.Atoi(rabbitMQPort)
		if err != nil {
			logger.Error("Invalid RABBITMQ_PORT value: %v; setting to default", err)
			intPort = 5672 // default RabbitMQ port
		}
		cfg.RabbitMQ.Port = intPort
	}
	if rabbitMQUser := os.Getenv("RABBITMQ_USER"); rabbitMQUser != "" {
		cfg.RabbitMQ.User = rabbitMQUser
	}
	if rabbitMQPass := os.Getenv("RABBITMQ_PASS"); rabbitMQPass != "" {
		cfg.RabbitMQ.Password = rabbitMQPass
	}
}
func setFrontendConfigByEnv(cfg *Config) {
	if frontendBaseURL := os.Getenv("FRONTEND_BASE_URL"); frontendBaseURL != "" {
		cfg.Frontend.BaseURL = frontendBaseURL
	}
	if confirmationEndpoint := os.Getenv("CONFIRMATION_ENDPOINT"); confirmationEndpoint != "" {
		cfg.Frontend.ConfirmationEndpoint = confirmationEndpoint
	}
}

func setGeneralConfigByEnv(cfg *Config) {
	if appName := os.Getenv("APP_NAME"); appName != "" {
		cfg.AppName = appName
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
}
