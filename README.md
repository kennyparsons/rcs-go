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

### Example `project.toml`

```toml
[ssh]
host = "my-server"

[exec]
dir = "/srv/app"
shell = false
timeout = "5m"

[env]
API_KEY = "default_key"
```
