package config

import "time"

type Config struct {
	Debug    bool           `mapstructure:"debug" yaml:"debug"`
	App      AppConfig      `mapstructure:"app" yaml:"app"`
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Auth     AuthConfig     `mapstructure:"auth" yaml:"auth" validate:"required"`
	Kafka    KafkaConfig    `mapstructure:"kafka" yaml:"kafka" validate:"required"`
}

type AppConfig struct {
	Name string `mapstructure:"name" yaml:"name"`
	Env  string `mapstructure:"env" yaml:"env"`
}

type ServerConfig struct {
	Port            int           `mapstructure:"port" yaml:"port"`
	Host            string        `mapstructure:"host" yaml:"host"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	User     string `mapstructure:"user" yaml:"user"`
	Password string `mapstructure:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" yaml:"dbname"`
	SSLMode  string `mapstructure:"sslmode" yaml:"sslmode"`
}

type AuthConfig struct {
	JWTSecret string        `mapstructure:"jwt_secret" yaml:"jwt_secret" validate:"required"`
	JWTExpiry time.Duration `mapstructure:"jwt_expiry" yaml:"jwt_expiry" validate:"required"`
}

type KafkaConfig struct {
	Brokers       []string `mapstructure:"brokers" yaml:"brokers" validate:"required"`
	ConsumerGroup string   `mapstructure:"consumer_group" yaml:"consumer_group" validate:"required"`
	Topics        Topics   `mapstructure:"topics" yaml:"topics" validate:"required"`
}

type TopicConfig struct {
	Name        string `mapstructure:"name" yaml:"name"`
	Concurrency int    `mapstructure:"concurrency" yaml:"concurrency"`
}

type Topics map[string]TopicConfig

func (t Topics) GetName(topic string) string {
	if t == nil {
		return topic
	}

	topicCfg := t[topic]
	if topicCfg.Name != "" {
		return topicCfg.Name
	}

	return topic
}
