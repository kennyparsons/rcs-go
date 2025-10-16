package exec

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
	"github.com/kennyparsons/rcs/internal/types"
)

// BuildRemoteCommand constructs the command string to be executed remotely.
func BuildRemoteCommand(opts *types.CLIOptions) string {
	var envs string
	if len(opts.Env) > 0 {
		var envParts []string
		for _, e := range opts.Env {
			envParts = append(envParts, fmt.Sprintf("export %s", e))
		}
		envs = strings.Join(envParts, " && ") + " && "
	}

	var cd string
	if opts.Dir != "" {
		cd = fmt.Sprintf("cd %s && ", opts.Dir)
	}

	command := strings.Join(opts.Command, " ")
	if opts.Shell {
		// In shell mode, wrap the command to be executed by bash.
		// This allows for pipes, redirection, etc.
		escapedCommand := strings.ReplaceAll(command, "'", "'\\''")
		return fmt.Sprintf("%s%sbash -lc '%s'", cd, envs, escapedCommand)
	}

	// In argv mode, the command is executed directly.
	return fmt.Sprintf("%s%s%s", cd, envs, command)
}

// ExitStatusOf returns the exit status code from an error returned by an SSH session.
func ExitStatusOf(err error) int {
	if err == nil {
		return 0
	}
	if exitErr, ok := err.(*ssh.ExitError); ok {
		return exitErr.ExitStatus()
	}
	// Return a non-zero status for other errors (e.g., connection failed).
	return 1
}
