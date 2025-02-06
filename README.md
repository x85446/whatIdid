# whatIdid

A personal activity tracking tool that logs and categorizes your digital interactions across various sources.

## Quick Start

1. Clone the repository
2. Copy the template config:
   ```bash
   cp config/dummy.yaml config/wid.yaml
   ```
3. Edit `config/wid.yaml` with your settings
4. Run the fetch command:
   ```bash
   go run main.go fetch
   ```

## Using Make Commands

The project includes several make targets to help with development:

```bash
# Build the application
make build

# Clean build artifacts and database
make clean

# Run all tests
make test

# Build and run the application
make run

# Initial setup (creates config/wid.yaml from template)
make setup

# Fetch new events
make fetch

# Install development dependencies
make dev-deps

# Run the linter
make lint
```

Common workflows:
1. First time setup: `make setup && make build`
2. Regular development: `make clean && make run`
3. Before committing: `make test && make lint`

## Configuration

The application looks for configuration in the following order:
1. `./config/wid.yaml` (recommended)
2. `$HOME/.config/whatidid/wid.yaml` (legacy)

See `config/dummy.yaml` for example configuration options.

## Database

Events are stored in SQLite database (`whatidid.db`). You can inspect the database using:

```bash
sqlite3 whatidid.db
```

Or using DB Browser for SQLite:
```bash
brew install --cask db-browser-for-sqlite
```

## Available Commands

- `fetch`: Collect new events from enabled plugins
- More commands coming soon...

## Development

See [requirements.md](docs/requirements.md) for detailed specifications and planned features.
