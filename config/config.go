package config

import (
	"strings"

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
	Mode        Mode
}

type Mode int

const (
	HTTPMode Mode = iota
	TestMode
	ProcessMode
)

func createMode(m string) Mode {
	v := strings.ToLower(m)
	switch v {
	case "http":
		return HTTPMode
	case "test":
		return TestMode
	default:
		return ProcessMode
	}
}

func (m Mode) String() string {
	switch m {
	case HTTPMode:
		return "HTTPMode"
	case TestMode:
		return "TestMode"
	default:
		return "ProcessMode"
	}
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
	conf.Mode = HTTPMode

	return &conf
}

func Get() *Config {
	return gConf
}
