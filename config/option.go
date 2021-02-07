package config

import "golang.org/x/xerrors"

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
