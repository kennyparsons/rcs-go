# rcs-go

A Go-based command-line tool for executing commands on remote hosts over SSH.

## Features

- **SSH Integration**: Uses your local SSH configuration (`~/.ssh/config`) for host aliases, keys, and settings.
- **Two Execution Modes**:
    - **argv mode (default)**: Executes commands directly without shell interpretation. This is the safer mode, preventing unintended shell expansion.
    - **shell mode (`--shell`)**: Executes commands within a remote login shell, allowing for pipes, redirection, and other shell features.
- **Configuration File**: Discovers and uses a `project.toml` or `project.json` file in the current or parent directories to set default values for host, directory, and other options.
- **I/O Streaming**: Streams remote stdout and stderr to the local terminal in real-time. Supports stdin passthrough for piping local data to remote commands.
- **PTY Allocation**: Can allocate a pseudo-terminal (`--pty`) for interactive, screen-based applications.
- **Environment Variables**: Injects environment variables into the remote session using the `--env` flag.

## Usage

The command requires a host (via the `-h` flag or a config file) and a command to execute, separated by `--`.

### Basic Example (argv mode)

```sh
./rcs-go -h my-server -- ls -la /var/www
```

### Shell Mode Example

```sh
./rcs-go -h my-server --shell -- 'grep "ERROR" /var/log/app.log | wc -l'
```

### Stdin Passthrough

```sh
cat local-data.txt | ./rcs-go -h my-server --stdin -- 'cat > /tmp/remote-data.txt'
```

## Configuration

Create a `project.toml` or `project.json` file in your project's root directory to define default settings. CLI flags will always override settings from the configuration file.

### Detailed `project.toml` Example

```toml
# SSH connection settings
[ssh]
# The default remote host or SSH alias to connect to.
host = "prod-server"

# The path to the SSH private key to use for authentication.
# Tilde (~) is supported for the home directory.
identity = "~/.ssh/id_ed25519_prod"

# The path to your SSH configuration file.
# Defaults to ~/.ssh/config if not specified.
ssh_config = "~/.ssh/config"

# Command execution settings
[exec]
# The default remote working directory to change into before running the command.
dir = "/srv/app/backend"

# If true, allocates a pseudo-terminal (PTY). Useful for interactive programs.
pty = false

# If true, runs the command inside a login shell (bash -lc '...').
# This allows for pipes, redirection, and sourcing of shell profiles.
shell = false

# If true, streams local stdin to the remote command.
stdin = false

# The maximum time to wait for the command to complete.
# Uses Go's duration format (e.g., "10s", "5m", "1h30m").
timeout = "10m"

# Logging settings
[log]
# If set, writes all stdout and stderr to this file.
path = "/var/log/rcs-go/deployment.log"

# If true and a log path is set, appends to the log file instead of overwriting it.
append = true

# If true, prepends a UTC timestamp to each line in the log file.
timestamps = true

# If true, splits stdout and stderr into separate files.
# For a path of "/tmp/log", this would create "/tmp/log.out" and "/tmp/log.err".
split = false

# Environment variables to set on the remote host before execution.
[env]
APP_ENV = "production"
DATABASE_URL = "prod_db_connection_string"
API_SECRET = "a-default-secret-from-config"
```

### Configuration Options

| Section | Key          | Type    | Description                                                                 |
|---------|--------------|---------|-----------------------------------------------------------------------------|
| `ssh`   | `host`       | String  | Default SSH host or alias.                                                  |
| `ssh`   | `identity`   | String  | Path to the SSH private key.                                                |
| `ssh`   | `ssh_config` | String  | Path to the SSH config file.                                                |
| `exec`  | `dir`        | String  | Default remote working directory.                                           |
| `exec`  | `pty`        | Boolean | Allocate a pseudo-terminal.                                                 |
| `exec`  | `shell`      | Boolean | Execute the command in a remote shell.                                      |
| `exec`  | `stdin`      | Boolean | Stream local stdin to the remote command.                                   |
| `exec`  | `timeout`    | String  | Command timeout (e.g., "30s", "5m").                                        |
| `log`   | `path`       | String  | Path to the output log file.                                                |
| `log`   | `append`     | Boolean | Append to the log file instead of overwriting.                              |
| `log`   | `timestamps` | Boolean | Add UTC timestamps to log entries.                                          |
| `log`   | `split`      | Boolean | Split stdout and stderr into separate log files (`.out` and `.err`).        |
| `env`   | `[key]`      | String  | Any key-value pair in this section is set as a remote environment variable. |
