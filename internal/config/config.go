package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kennyparsons/rcs/internal/types"
)

// Load finds and parses a configuration file (project.toml or project.json)
// and merges it with the provided CLI options. CLI options take precedence.
func Load(opts *types.CLIOptions) (*types.CLIOptions, error) {
	cfg, err := findAndParseConfig(opts.ConfigPath)
	if err != nil {
		// It's okay if the config file doesn't exist.
		if os.IsNotExist(err) {
			return opts, nil
		}
		return nil, err
	}

	// Merge config file values into opts, only if the flag was not set.
	if opts.Host == "" {
		opts.Host = cfg.SSH.Host
	}
	if opts.Dir == "" {
		opts.Dir = cfg.Exec.Dir
	}
	if !opts.Shell && cfg.Exec.Shell {
		opts.Shell = cfg.Exec.Shell
	}
	if !opts.Stdin && cfg.Exec.Stdin {
		opts.Stdin = cfg.Exec.Stdin
	}
	if !opts.Pty && cfg.Exec.Pty {
		opts.Pty = cfg.Exec.Pty
	}
	if opts.LogPath == "" {
		opts.LogPath = cfg.Log.Path
	}
	if !opts.Append && cfg.Log.Append {
		opts.Append = cfg.Log.Append
	}
	if !opts.Timestamps && cfg.Log.Timestamps {
		opts.Timestamps = cfg.Log.Timestamps
	}
	if !opts.Split && cfg.Log.Split {
		opts.Split = cfg.Log.Split
	}
	if opts.Identity == "" {
		opts.Identity = cfg.SSH.Identity
	}
	if opts.SSHConfigPath == "" {
		opts.SSHConfigPath = cfg.SSH.SSHConfig
	}
	if opts.Timeout == 0 && cfg.Exec.Timeout != "" {
		if d, err := time.ParseDuration(cfg.Exec.Timeout); err == nil {
			opts.Timeout = d
		}
	}

	// Prepend env vars from config file to CLI env vars.
	if len(cfg.Env) > 0 {
		var configEnvs []string
		for k, v := range cfg.Env {
			configEnvs = append(configEnvs, k+"="+v)
		}
		opts.Env = append(configEnvs, opts.Env...)
	}

	return opts, nil
}

func findAndParseConfig(configPath string) (*types.Config, error) {
	var path string
	var err error

	if configPath != "" {
		path = configPath
	} else {
		path, err = findConfigFile()
		if err != nil {
			return nil, err
		}
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &types.Config{}
	if filepath.Ext(path) == ".toml" {
		err = toml.Unmarshal(file, cfg)
	} else {
		err = json.Unmarshal(file, cfg)
	}

	return cfg, err
}

func findConfigFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		for _, name := range []string{"project.toml", "project.json"} {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", os.ErrNotExist
}
