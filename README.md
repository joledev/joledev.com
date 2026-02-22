# joledev.com

Landing page, portfolio, blog, interactive quoter, and appointment scheduler.

## Stack

- **Frontend:** Astro 5 + Svelte 5 + Tailwind CSS 4
- **Backend:** Go microservices (Chi router)
- **Database:** SQLite + Litestream
- **Deploy:** Docker Compose + Traefik + GitHub Actions

## Development

### Frontend

```bash
cd apps/web
pnpm install
pnpm dev
```

### Go APIs

```bash
cd apps/api-quoter
go run main.go

cd apps/api-scheduler
go run main.go
```

### Docker (production)

```bash
cp .env.example .env  # Edit with real values
docker compose up -d
```

## Structure

```
apps/
├── web/              # Astro static site + Svelte islands
├── api-quoter/       # Quote calculation + email API
├── api-scheduler/    # Appointment booking API
└── status/           # Gatus monitoring config
```
