package config

import (
	"golang.org/x/xerrors"
)

type Config struct {
	Port        int
	Verbose     bool
	Schema      string
	Server      string
	Bin         string
	Args        []string
	IgnoreFiles []string
}

var gConf *Config

func init() {
	gConf = defaultConfig()
}

func Set(opts []Option) error {
	for _, opt := range opts {
		err := opt(gConf)
		if err != nil {
			return xerrors.Errorf("option setting error: %w", err)
		}
	}
	return nil
}

func defaultConfig() *Config {

	conf := Config{}
	conf.Port = 8080
	conf.Verbose = false
	conf.Schema = "http"
	conf.Server = "localhost"
	conf.Args = nil
	conf.IgnoreFiles = nil
	conf.Bin = "skewer-bin"

	return &conf
}

func Get() *Config {
	return gConf
}
