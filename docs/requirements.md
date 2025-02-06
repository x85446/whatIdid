# whatIdid - Features & Requirements Document

**Version:** 1.0

**Author:** Travis McCollum & AI

## 1. Overview

whatIdid is a Golang-based personal activity tracking tool that records and categorizes a user’s digital interactions across multiple applications. The tool provides a structured database of minute-by-minute events and enables AI-powered report generation.

### Core Objectives:

- **Automated event tracking:** Collects data from multiple sources (emails, browser history, coding activity, etc.).
- **Plugin-based architecture:** Each data source is an independent plugin.
- **AI-assisted reporting:** AI helps filter and categorize events into meaningful reports.
- **Flexible querying:** Users can filter activity logs based on any criteria (e.g., Kubernetes-related work, total time spent watching YouTube, etc.).
- **Local storage only:** Uses an internal, file-based database with no cloud dependencies.
- **Lightweight CLI:** Runs on demand (not a daemon), with optional Web UI for browsing logs.

## 2. Core Features

### 2.1. Event Collection

Plugins extract activity logs from various sources.

Each event includes:

- **Timestamp (Primary Index):** Standardized format for all events.
- **Source:** The plugin that generated the event (e.g., Chrome, Git, Outlook).
- **Event Type:** The nature of the event (e.g., email-sent, git-commit, browser-history).
- **Metadata:** Additional contextual details (e.g., email subject, URL, commit message).

Overlapping events are logged separately.

No idle time tracking – only definitive actions are recorded.

**Example Logs:**

- `2024-10-10T12:30:00Z, youtube, "Frankenstein’s Bride"`
- `2024-10-10T12:31:00Z, email-sent, "incus architecture discussion"`
- `2024-10-10T12:32:00Z, git-commit, "Refactored incus deployment"`

### 2.2. CLI Commands


**Fetching Data**

```bash
./whatIdid --config file fetch
```

- `Collects new events from all enabled plugins.`
- `plugin is called from a plugin controller with the start and stop dates.`
  - `plugin controller maintains a database of all dates ever fetched per controller`
  - `plugin controller tells plugin via an object what dates to skip in the fetch`
- `Dates are tracked to the day.`
- `events are tracked to the timestamp in RFC3339 format.`
- `plugins are responsible to normaize the timestamps from the format that they fetch data from, and place in the database in RFC3339 format.`

```bash
./whatIdid fetch --stop 2024-09-01 --start: 2025-01-01
```

- `Generating Reports (AI-Assisted)`

```bash
./whatIdid report
```
- `Starts an interactive session where the AI helps filter and summarize logs.`
- `Users specify criteria (e.g., "Show me all api-related activity between Dec 1 and Jan 1, 2025").`
- `The AI selects relevant events and formats them into a structured report.`
- `Launching Web UI (Future Feature)`

```bash 
./whatIdid web
```

- `Starts a lightweight local web interface for browsing and querying logs.`

- `REST API backend for UI and CLI integration.`

## 2.3. Plugin System

Plugins are configured via YAML (default location: `~/.config/whatIdid/wid.yaml`).

- `Users can enable/disable plugins in the config file.`
- ``
- `Plugins must conform to a standardized event format.`
- ``
- `New plugins can be added dynamically.`

**Example wid.yaml:**
```yaml
plugins:
  outlook:
    enabled: true
  chrome:
    enabled: true
    history_path: "/home/user/.config/google-chrome/Default/History"
  slack:
    enabled: false
  git:
    enabled: true
```

## 2.4. Database & Storage

- `SQLite or BoltDB (file-based) as the primary database.`
- `Primary index on timestamp for fast lookups.`
- `Normalized event table:`

```sql
CREATE TABLE events (
    timestamp TEXT PRIMARY KEY,
    source TEXT NOT NULL,
    event_type TEXT NOT NULL,
    metadata TEXT
);
```

- `Plugins may have additional tables if needed.`
- `File-based snapshots for backups and versioning.`

# 2.5. AI Integration (Future Feature)

- `AI will analyze and categorize events without modifying raw data.`
- `AI helps users generate meaningful reports through an interactive session.`
- `Categorization will be stored separately, allowing AI and users to refine classifications collaboratively.`

# 2.6. Security & Privacy

- `Runs on demand (not a background daemon).\`
- `Local storage only – no cloud sync.`
- `Users can define a deletion policy in wid.yaml:`

```yaml
retention:
  delete_after_days: 90  # Delete records older than 90 days
```
# Plugin API Design (plugin.go)

```go
package plugin

type Event struct {
    Timestamp string `json:"timestamp"`
    Source    string `json:"source"`
    EventType string `json:"event_type"`
    Metadata  string `json:"metadata"`
}

type Plugin interface {
    Name() string
    FetchEvents(lastFetchTime string) ([]Event, error)
}
```
# Next Steps

- `Implement the core CLI commands (fetch, report, web).`
- `Develop initial plugins (Chrome history, Git commits, Outlook emails).`
- `Build the database layer and ensure efficient querying`
- `Design the REST API for the Web UI.`
- `Does this updated document fully capture the implementation details? 🚀`

