# gstest - Source Engine Game Server Tester

A CLI tool for testing Source engine game servers using the A2S protocol. Designed for use in CI/CD pipelines, container health checks, and manual server validation.

## Features

- **Connectivity check**: Verifies server responds to A2S_INFO queries
- **Map loaded check**: Ensures a map is loaded and not in a loading state
- **Player slots check**: Confirms players can join (max_players > 0 and not full)
- **JSON output**: Structured output for scripting and automation
- **Exit codes**: Specific codes for different failure types

## Installation

### From source

```bash
go install gameserver-testing/cmd/gstest@latest
```

### Using Docker

```bash
docker pull gstest:latest
```

### Build from source

```bash
git clone <repository>
cd gameserver-testing
make build
```

## Usage

```bash
# Basic connectivity test
gstest 192.168.1.100

# Specify port and timeout
gstest 192.168.1.100 --port 27016 --timeout 10s

# Run specific checks only
gstest 192.168.1.100 --checks connectivity,maploaded

# JSON output for scripting
gstest 192.168.1.100 --json

# Docker usage
docker run --rm gstest:latest 192.168.1.100 --port 27015 --json
```

## CLI Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | 27015 | Server port |
| `--timeout` | `-t` | 5s | Query timeout |
| `--checks` | `-c` | all | Comma-separated list of checks |
| `--json` | `-j` | false | Output results as JSON |
| `--verbose` | `-v` | false | Verbose output |
| `--version` | | | Show version information |
| `--help` | `-h` | | Show help |

## Available Checks

| Check | Description |
|-------|-------------|
| `connectivity` | Server responds to A2S_INFO query |
| `maploaded` | Map field is populated (not empty or loading) |
| `playerslots` | max_players > 0 and server is not full |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | All checks passed |
| 1 | Connectivity check failed |
| 2 | Map not loaded check failed |
| 3 | Player slots check failed (server full) |
| 10 | Configuration error |
| 99 | Unknown error |

## JSON Output Example

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "server": {
    "address": "192.168.1.100:27015",
    "name": "My CS2 Server",
    "map": "de_dust2",
    "game": "Counter-Strike 2",
    "players": 12,
    "max_players": 24,
    "bots": 0,
    "server_type": "Dedicated",
    "vac": true
  },
  "results": [
    {
      "name": "connectivity",
      "passed": true,
      "message": "Server responded: My CS2 Server"
    },
    {
      "name": "maploaded",
      "passed": true,
      "message": "Map loaded: de_dust2"
    },
    {
      "name": "playerslots",
      "passed": true,
      "message": "Player slots available: 12/24 (free: 12)"
    }
  ],
  "all_passed": true,
  "exit_code": 0
}
```

## Docker Compose Example

```yaml
services:
  health-check:
    image: gstest:latest
    command: ["gameserver", "--port", "27015", "--json"]
    network_mode: host
```

## Development

```bash
# Build
make build

# Run tests
make test

# Build Docker image
make docker

# Run linter
make lint
```

## License

MIT
