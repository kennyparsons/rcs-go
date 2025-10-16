package types

import "time"

// CLIOptions holds all the command-line flags.
type CLIOptions struct {
	Host          string
	Dir           string
	Shell         bool
	Stdin         bool
	Pty           bool
	LogPath       string
	Append        bool
	Timestamps    bool
	Split         bool
	Identity      string
	SSHConfigPath string
	ConfigPath    string
	Timeout       time.Duration
	Env           []string
	Command       []string
}

// Config represents the structure of the project.toml or project.json file.
type Config struct {
	SSH  SSHConfig  `json:"ssh" toml:"ssh"`
	Exec ExecConfig `json:"exec" toml:"exec"`
	Log  LogConfig  `json:"log" toml:"log"`
	Env  EnvConfig  `json:"env" toml:"env"`
}

// SSHConfig holds SSH-related settings.
type SSHConfig struct {
	Host       string `json:"host" toml:"host"`
	Identity   string `json:"identity" toml:"identity"`
	SSHConfig  string `json:"ssh_config" toml:"ssh_config"`
}

// ExecConfig holds execution-related settings.
type ExecConfig struct {
	Dir     string `json:"dir" toml:"dir"`
	Pty     bool   `json:"pty" toml:"pty"`
	Shell   bool   `json:"shell" toml:"shell"`
	Stdin   bool   `json:"stdin" toml:"stdin"`
	Timeout string `json:"timeout" toml:"timeout"`
}

// LogConfig holds logging-related settings.
type LogConfig struct {
	Path       string `json:"path" toml:"path"`
	Append     bool   `json:"append" toml:"append"`
	Timestamps bool   `json:"timestamps" toml:"timestamps"`
	Split      bool   `json:"split" toml:"split"`
}

// EnvConfig is a map for environment variables.
type EnvConfig map[string]string
