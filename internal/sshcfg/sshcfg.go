package sshcfg

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/kevinburke/ssh_config"
)

// Target represents a resolved SSH target with all necessary connection details.
type Target struct {
	Host         string
	User         string
	Port         int
	IdentityFile string
}

// ResolveHost resolves an SSH host alias using the user's SSH config.
func ResolveHost(host, sshConfigPath string) (*Target, error) {
	if sshConfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		sshConfigPath = filepath.Join(home, ".ssh", "config")
	}

	f, err := os.Open(sshConfigPath)
	if err != nil {
		// A missing SSH config is not a fatal error.
		if os.IsNotExist(err) {
			return &Target{Host: host, Port: 22}, nil
		}
		return nil, err
	}
	defer f.Close()

	cfg, err := ssh_config.Decode(f)
	if err != nil {
		return nil, err
	}

	hostname, err := cfg.Get(host, "HostName")
	if err != nil || hostname == "" {
		hostname = host
	}

	user, _ := cfg.Get(host, "User")

	portStr, _ := cfg.Get(host, "Port")
	port := 22
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	identityFile, _ := cfg.Get(host, "IdentityFile")

	return &Target{
		Host:         hostname,
		User:         user,
		Port:         port,
		IdentityFile: identityFile,
	}, nil
}
