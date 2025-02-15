package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServiceName  string              `yaml:"serviceName" env-required:"true"`
	SrvDiscovery *SrvDiscoveryConfig `yaml:"srv_discovery"`
	Server       *ServerConfig       `yaml:"server"`
	DB           *DBConfig           `yaml:"db"`
	Redis        *RedisConfig        `yaml:"redis"`
	Jaeger       *JaegerConfig       `yaml:"jaeger"`
}

type SrvDiscoveryConfig struct {
	Scheme string `yaml:"scheme" env-default:"http"`
	Host   string `yaml:"host" env-default:"localhost"`
	Port   int    `yaml:"port" env-default:"50030"`
}

type ServerConfig struct {
	Port   int    `yaml:"port" env-required:"true"`
	Mode   string `yaml:"mode" env-default:"dev"`
	Scheme string `yaml:"scheme" env-default:"http"`
	Domain string `yaml:"domain" env-default:"localhost"`
}

type DBConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Database string `yaml:"database" env-default:"db"`
}

type RedisConfig struct {
	Addr string `yaml:"addr" env-default:"localhost:6379"`
	Pass string `yaml:"pass" env-default:""`
}

type JaegerConfig struct {
	Sampler struct {
		Type  string `yaml:"type"`
		Param int    `yaml:"param"`
	} `yaml:"sampler"`
	Reporter struct {
		LogSpans           bool   `yaml:"LogSpans"`
		LocalAgentHostPort string `yaml:"LocalAgentHostPort"`
	} `yaml:"reporter"`
}

func MustLoad(configPath string) *Config {
	var conf Config

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	if err = yaml.Unmarshal(data, &conf); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	return &conf
}
