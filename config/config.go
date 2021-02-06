package config

import (
	"runtime"

	"golang.org/x/xerrors"
)

type Config struct {
	Port    int
	AppPort int
	Verbose bool
	Schema  string
	Server  string
	Bin     string
	Args    []string
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
	conf.Port = 3000
	conf.AppPort = 8080
	conf.Verbose = false
	conf.Schema = "http"
	conf.Server = "localhost"
	conf.Args = nil

	exe := ""
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}
	conf.Bin = ".skewer-work-bin" + exe

	return &conf
}

func Get() *Config {
	return gConf
}
