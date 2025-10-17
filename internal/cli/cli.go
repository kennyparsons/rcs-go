package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/kennyparsons/rcs/internal/types"
)

// stringSlice is a custom type to handle repeatable string flags.
type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Parse parses the command-line arguments and returns the configured CLIOptions.
// It also returns a boolean indicating if the help flag was requested.
func Parse() (*types.CLIOptions, bool, error) {
	opts := &types.CLIOptions{}
	var envs stringSlice
	var help bool

	fs := flag.NewFlagSet("rcs-go", flag.ContinueOnError)

	fs.StringVar(&opts.Host, "h", "", "SSH host or alias")
	fs.StringVar(&opts.Host, "host", "", "SSH host or alias")
	fs.StringVar(&opts.Dir, "dir", "", "Remote working directory")
	fs.BoolVar(&opts.Shell, "shell", false, "Execute command in a shell")
	fs.BoolVar(&opts.Stdin, "stdin", false, "Stream local stdin to the remote process")
	fs.BoolVar(&opts.Pty, "pty", false, "Allocate a PTY for the remote session")
	fs.StringVar(&opts.LogPath, "log", "", "Log output to a file")
	fs.BoolVar(&opts.Append, "append", false, "Append to the log file if it exists")
	fs.BoolVar(&opts.Timestamps, "timestamps", false, "Add timestamps to log entries")
	fs.BoolVar(&opts.Split, "split", false, "Split stdout and stderr into separate log files")
	fs.StringVar(&opts.Identity, "identity", "", "Path to SSH private key")
	fs.StringVar(&opts.SSHConfigPath, "ssh-config", "", "Path to SSH config file")
	fs.StringVar(&opts.ConfigPath, "config", "", "Path to project config file (JSON or TOML)")
	fs.DurationVar(&opts.Timeout, "timeout", 0, "Command timeout")
	fs.Var(&envs, "env", "Set environment variables remotely (e.g., KEY=VAL)")
	fs.BoolVar(&help, "help", false, "Show help message")

	// Find the command separator "--"
	args := os.Args[1:]
	separatorIndex := -1
	for i, arg := range args {
		if arg == "--" {
			separatorIndex = i
			break
		}
	}

	var commandArgs []string
	if separatorIndex != -1 {
		// Parse flags before the separator
		if err := fs.Parse(args[:separatorIndex]); err != nil {
			return nil, false, err
		}
		// The rest are command arguments
		if separatorIndex+1 < len(args) {
			commandArgs = args[separatorIndex+1:]
		}
	} else {
		// No separator, parse all args as flags
		if err := fs.Parse(args); err != nil {
			return nil, false, err
		}
		// Command is the remaining non-flag arguments
		commandArgs = fs.Args()
	}

	if help {
		fs.Usage()
		return nil, true, nil
	}

	if len(commandArgs) == 0 {
		return nil, false, fmt.Errorf("error: command is required")
	}

	opts.Env = envs
	opts.Command = commandArgs

	return opts, false, nil
}
