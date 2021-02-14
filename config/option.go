package config

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

type Option func(*Config) error

func SetArgs(args []string) Option {
	return func(c *Config) error {
		if args == nil || len(args) <= 0 {
			return xerrors.Errorf("build files required.")
		}
		c.Args = args
		return nil
	}
}

func SetVerbose(v bool) Option {
	return func(c *Config) error {
		c.Verbose = v
		return nil
	}
}

func SetMode(m string, p int, f bool) Option {
	return func(c *Config) error {
		mode := createMode(m)
		c.Mode = mode
		if mode == ProcessMode {
			return fmt.Errorf("%s Mode is not implemented.", mode)
		}

		if f {
			//HTTP mode only
			portBuf := os.Getenv("PORT")
			if portBuf == "" {
				return fmt.Errorf(`if "e" is specified as an argument,it must be specified in the "PORT" environment variable.`)
			}
			port, err := strconv.Atoi(portBuf)
			if err != nil {
				return fmt.Errorf(`the "PORT" environment variable is not number.[%s]`, portBuf)
			}
			c.Port = port
		} else {
			c.Port = p
		}
		return nil
	}
}

func SetIgnoreFiles(f string) Option {
	return func(c *Config) error {
		files := strings.Split(f, "|")
		c.IgnoreFiles = files
		return nil
	}
}

func SetBin(n string) Option {
	return func(c *Config) error {
		exe := ""
		if runtime.GOOS == "windows" {
			exe = ".exe"
		}
		c.Bin = n + exe
		return nil
	}
}
