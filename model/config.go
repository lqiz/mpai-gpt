package model

import (
	"github.com/pelletier/go-toml"
	"os"
)

// Config example config
type Config struct {
	DB struct {
		DataSource string `toml:"datasource"`
	} `toml:"db"`
	Listen struct {
		Port string `toml:"port"`
	} `toml:"listen"`
	Quota struct {
		Max int64 `toml:"max"`
	} `toml:"quota"`
	RemoteProxy struct {
		Url string `toml:"url"`
	} `toml:"remoteProxy"`
	Redis struct {
		Host        string `toml:"host"`
		Password    string `toml:"password"`
		Database    int    `toml:"database"`
		MaxActive   int    `toml:"maxActive"`
		MaxIdle     int    `toml:"maxIdle"`
		IdleTimeout int    `toml:"idleTimeout"`
	} `toml:"redis"`
	*OfficialAccountConfig `toml:"officialAccountConfig"`
}

// OfficialAccountConfig 公众号相关配置
type OfficialAccountConfig struct {
	AppID          string `toml:"appID"`
	AppSecret      string `toml:"appSecret"`
	Token          string `toml:"token"`
	EncodingAESKey string `toml:"encodingAESKey"`
}

// GetConfig 获取配置
func GetConfig(cfgFile *string) *Config {
	bytes, err := os.ReadFile(*cfgFile)
	if err != nil {
		panic(err)
	}

	cfgData := &Config{}
	err = toml.Unmarshal(bytes, cfgData)
	if err != nil {
		panic(err)
	}
	return cfgData
}
