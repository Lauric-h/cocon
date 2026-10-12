# Cocon

Cocon is a lightweight household dashboard for planning everyday life in one place. It combines weekly meal planning, grocery tracking, tasks, calendar events, and a downloadable weekly summary.

The application is intentionally small and practical: a server-rendered Go web app backed by SQLite, with responsive HTML/CSS and HTMX for fast partial updates without a heavy frontend framework.

## Features

- Weekly dashboard with meals, tasks, and calendar events
- Meal planning for breakfast, lunch, and dinner
- Grocery list with quantities and completion tracking
- Planned tasks and backlog tasks with categories and due dates
- Google Calendar integration through an ICS feed
- Weekly summary exported as PDF
- Responsive desktop and mobile layouts
- Persistent local SQLite storage

## Tech stack

- **Go 1.26** — HTTP server, application logic, HTML templates
- **SQLite** — local persistence with automatic schema initialization
- **HTML/CSS + HTMX** — responsive server-rendered interface and partial updates
- **Docker** — reproducible production-style deployment
- **GitHub Actions** — multi-architecture Docker image publishing

## Run with Docker Compose

Requirements: Docker with the Compose plugin.

```bash
docker compose up -d
```

Open <http://localhost:8081>.

The Compose setup:

- uses the published image `ghcr.io/lauric-h/cocon:latest`
- exposes the application on port `8081`
- persists SQLite data in `./data`
- reads the optional `GOOGLE_CALENDAR_ICS_URL` variable from the environment

Stop the stack:

```bash
docker compose down
```

To enable calendar events, export an ICS URL before starting the stack:

```bash
GOOGLE_CALENDAR_ICS_URL="https://example.com/calendar.ics" docker compose up -d
```

## Run locally

Requirements: Go 1.26 or newer.

```bash
mkdir -p data
PORT=8081 DB_PATH=./data/cocon.db go run ./cmd/web
```

Open <http://localhost:8081>.

Environment variables:

| Variable                  | Default                  | Purpose                          |
| ------------------------- | ------------------------ | -------------------------------- |
| `PORT`                    | `8081`                   | HTTP server port                 |
| `DB_PATH`                 | `./data/cocon.db`        | SQLite database path             |
| `TZ`                      | `Europe/Paris` in Docker | Application timezone             |
| `GOOGLE_CALENDAR_ICS_URL` | empty                    | Optional read-only calendar feed |

The application also loads variables from a local `.env` file. Keep credentials and private calendar URLs out of version control.

## Project structure

```text
cmd/web/                 Application entrypoint
internal/handlers/       HTTP routes and view models
internal/store/          SQLite persistence
internal/calendar/       ICS calendar client and cache
internal/pdfexport/      Weekly PDF generation
ui/html/                 Go HTML templates
ui/static/               CSS assets
data/                    Local SQLite database files
```

## Engineering highlights

- Clear separation between HTTP handlers, domain models, storage, integrations, and presentation.
- Standard-library HTTP routing with no frontend build step.
- SQLite schema created at startup, keeping local development simple.
- Docker builds a small multi-stage image and supports `linux/amd64` and `linux/arm64`.
- Calendar events are cached in memory to avoid fetching the ICS feed on every request.
