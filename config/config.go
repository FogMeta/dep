package config

import (
	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Port     int      `toml:"port"`
	DataBase Database `toml:"db"`
	Redis    Redis    `toml:"redis"`
	Email    Email    `toml:"email"`
	Lagrange Lagrange `toml:"lagrange"`
	Log      Log      `toml:"log"`
}

type Database struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	Debug    bool   `toml:"debug"`
}

type Redis struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type Email struct {
	Email           string `toml:"email"`
	Password        string `toml:"password"`
	Host            string `toml:"host"`
	Port            int    `toml:"port"`
	Subject         string `toml:"subject"`
	ContentTemplate string `toml:"content_template"`
}

type Lagrange struct {
	Host    string `toml:"host"`
	SDKHost string `toml:"sdk_host"`
}

type Log struct {
	Level int `toml:"level"`
}

var conf Configuration

const (
	confPath = "config/config.toml"
)

func Init() error {
	_, err := toml.DecodeFile(confPath, &conf)
	if err != nil {
		panic(err)
	}
	return nil
}

func Conf() *Configuration {
	return &conf
}
