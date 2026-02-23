# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**joledev.com** — Freelance developer portfolio, blog, interactive quoter, and appointment scheduler for Joel Ernesto López Verdugo (Ensenada, Baja California, México). The primary language of the site content is Spanish; English is secondary.

## Tech Stack

- **Frontend:** Astro 5 + Svelte 5 (islands architecture) + Tailwind CSS 4
- **3D Effects:** Three.js (direct, custom shaders for particles)
- **Animations:** Motion (Svelte) + CSS animations + Anime.js
- **Backend:** Go 1.23 with Chi router (two separate microservices)
- **Database:** SQLite (WAL mode) + Litestream (S3 backups)
- **Blog:** Markdown/MDX via Astro Content Collections
- **Email:** SMTP via Hostinger (implicit TLS, port 465) — single account: `contacto@joledev.com`
- **CAPTCHA:** Cloudflare Turnstile (quoter + scheduler forms)
- **i18n:** Astro i18n routing (es primary, en secondary)
- **Deploy:** GitHub Actions → GHCR → K3s (Traefik ingress + cert-manager TLS)
- **Monitoring:** Gatus (status page at status.joledev.com)

## Build & Development Commands

```bash
# Frontend (from apps/web/)
pnpm install            # Install dependencies
pnpm dev                # Dev server (http://localhost:4321)
pnpm build              # Production build → dist/

# Go APIs (from apps/api-quoter/ or apps/api-scheduler/)
go run main.go          # Run locally
go build -o bin/server  # Build binary
go test ./...           # Run all tests

# Docker — dev mode (from repo root)
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d api-quoter api-scheduler
docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build api-scheduler  # rebuild

# Docker — stop everything
docker compose -f docker-compose.yml -f docker-compose.dev.yml down

# Docker — production (from repo root)
docker compose up -d    # Start all services (traefik, nginx, apis, gatus)
docker compose build    # Rebuild images
```

## Architecture

### Monorepo with 4 apps under `apps/`:

- **`apps/web/`** — Astro static site with Svelte islands. Generates pure HTML/CSS/JS served by Nginx. Interactive components (quoter, scheduler, Three.js hero) are Svelte islands hydrated client-side.
- **`apps/api-quoter/`** — Go microservice (port 8081). Handles quote calculation and email sending via SMTP. Endpoint: `POST /quotes`.
- **`apps/api-scheduler/`** — Go microservice (port 8082). Manages appointment booking with auto-generated availability (Mon-Fri 9:00-16:00 America/Tijuana, 30min slots, 2h buffer). SQLite for persistence. Admin confirm/reject via email token links.
- **`apps/status/`** — Gatus config for monitoring all services.

### Scheduler System (recently redesigned)

The scheduler computes availability on-the-fly (no manual slot creation):
- **Availability:** Mon-Fri 9:00-15:30, every 30 min, America/Tijuana timezone
- **Buffer:** 2 hours between booking start times (max ~4 bookings/day)
- **Booking flow:** Client creates booking (status: `pending`) → admin gets email with confirm/reject links → client notified of result
- **One active booking per email** (pending or confirmed, date >= today)
- **Admin endpoints** protected with Basic Auth (`SCHEDULER_ADMIN_PASSWORD`)
- **Token-based confirm/reject:** `GET /scheduler/bookings/confirm?token=` and `GET /scheduler/bookings/reject?token=` (return HTML pages)
- **Key files:**
  - `services/availability.go` — Slot computation algorithm
  - `services/tokens.go` — Crypto/rand token generation
  - `handlers/bookings.go` — CreateBooking, ConfirmBooking, RejectBooking, GetAdminBookings, CancelBooking
  - `handlers/slots.go` — GetAvailableSlots (public)
- **Database:** Single `bookings` table with `date`, `start_time`, `end_time`, `confirm_token`, `reject_token`, `status`
- **Important:** Nullable columns (client_phone, client_company, client_address, notes) use `COALESCE(column, '')` in SELECT queries to avoid Go scan errors

### Production (K3s on VPS `69.62.68.130`):
- **Orchestrator:** K3s (lightweight Kubernetes) with Traefik ingress controller
- **TLS:** cert-manager with `letsencrypt-prod` ClusterIssuer
- **Registry:** `ghcr.io/joledev/` (joledev-web, joledev-api-quoter, joledev-api-scheduler)
- **Namespace:** `joledev`
- **Storage:** `local-path` StorageClass for SQLite PVC (api-scheduler)
- **Manifests:** `k8s/` directory (namespace, deployments, services, ingress, gatus configmap)
- **Secrets:** `joledev-secrets` (SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS, SMTP_FROM, CONTACT_EMAIL, SCHEDULER_ADMIN_PASSWORD, TURNSTILE_SECRET_KEY) + `ghcr-secret`

### Routing (Traefik Ingress):
- `joledev.com` → web (nginx serving Astro static files)
- `api.joledev.com/quotes` → api-quoter:8081
- `api.joledev.com/scheduler` → api-scheduler:8082
- `api.joledev.com/health` → api-quoter:8081
- `status.joledev.com` → gatus:8080

### Frontend Key Directories (`apps/web/src/`):
- `components/quoter/` — Multi-step quoter (7 steps: project type → features → business size → current state → timeline → currency → result + contact form)
- `components/scheduler/` — Booking flow (Scheduler.svelte) + admin panel (AdminScheduler.svelte)
- `components/hero/` — Three.js icosahedron with glass material, particles, orbiting rings. Theme-reactive (MutationObserver on data-theme)
- `content/blog/{es,en}/` — Blog posts in Markdown with frontmatter
- `content/projects/{es,en}/` — Project case studies (12 projects, both languages)
- `pages/{es,en}/` — Astro i18n routing (root index.astro redirects by browser locale)
- `i18n/` — Translation JSON files (es.json, en.json)
- `lib/quoter-config.ts` — Editable pricing/service configuration (prices as base + multipliers)
- `styles/global.css` — Design system: CSS custom properties for both themes, Tailwind @theme integration

### Environment Variables
See `.env.example`. Key vars:
- `SMTP_HOST` / `SMTP_PORT` / `SMTP_USER` / `SMTP_PASS` / `SMTP_FROM` — SMTP credentials (Hostinger, implicit TLS on port 465)
- `CONTACT_EMAIL` — Recipient for quoter and scheduler notifications (`contacto@joledev.com`)
- `SCHEDULER_ADMIN_PASSWORD` — Basic Auth password for admin scheduler routes
- `API_BASE_URL` — Used by scheduler to build confirm/reject links in emails (prod: `https://api.joledev.com`, dev: `http://localhost:8082`)
- `DOMAIN` — Used by Traefik for routing rules
- `PUBLIC_TURNSTILE_SITE_KEY` / `TURNSTILE_SECRET_KEY` — Cloudflare Turnstile CAPTCHA keys (test keys in `.env.example` always pass)
- **Important:** Only one email account exists: `contacto@joledev.com`. Both `SMTP_FROM` and `CONTACT_EMAIL` use this address. Do not reference `noreply@` or other addresses.

## CI/CD Pipeline

- **CI (`ci.yml`):** Runs on PRs and pushes to non-main branches. Builds Astro frontend, builds Go APIs (CGO_ENABLED=1), runs Go tests.
- **Deploy (`deploy.yml`):** Runs on push to `main`. Three stages: CI → build & push images to GHCR → SSH to VPS and `kubectl rollout restart`.
- **Manual deploy (`scripts/deploy.sh`):** Restarts all K3s deployments in the `joledev` namespace.
- **Branch strategy:** Single `main` branch → production (joledev.com). No staging environment.
- **Required GitHub Secrets:** `VPS_HOST`, `VPS_PORT`, `VPS_USER`, `VPS_SSH_KEY`
- **GHCR auth:** Deploy workflow uses `GITHUB_TOKEN` (automatic) with `packages: write` permission.

## Security

- **Security headers:** Both APIs set `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy: strict-origin-when-cross-origin`
- **Rate limiting:** In-memory per-IP rate limiter on all public endpoints (quoter: 5/hr, scheduler bookings: 10/hr, slots: 60/hr)
- **Body size limits:** 64KB max on submissions, 4KB on admin actions
- **Input validation:** Email format regex, field length limits (name: 200, email: 254, phone: 30, company: 200, address: 500, notes: 2000), date (YYYY-MM-DD) and time (HH:MM) format validation
- **Admin auth:** Basic Auth with `subtle.ConstantTimeCompare` (timing-safe)
- **CORS:** Configurable via `CORS_ORIGIN` env var, defaults to `https://joledev.com`
- **CAPTCHA:** Cloudflare Turnstile on quoter contact form and scheduler booking form. Backend verifies token via `https://challenges.cloudflare.com/turnstile/v0/siteverify`. Test keys (always pass) used in dev.
- **Never commit:** `.env` files, SMTP credentials, API keys. Only `.env.example` with placeholders.

## Testing

```bash
# Run all tests (from each API directory)
CGO_ENABLED=1 go test ./...

# Verbose output
CGO_ENABLED=1 go test ./... -v
```

- **Quoter tests (6):** Valid request, missing email/name, empty project types, field length validation
- **Scheduler handler tests (15):** Booking CRUD, confirm/reject tokens, duplicate email, buffer enforcement, date/time format validation, meeting type validation, field length validation, slots date format
- **Scheduler service tests (5):** timeToMinutes/minutesToTime, weekday-only slots, slot count, 2h buffer blocking

## Code Conventions

- TypeScript strict mode in frontend
- Svelte 5 runes: `$state`, `$derived`, `$effect`, `$props`
- Svelte components: scoped styles, separate logic where possible
- Go: explicit error handling, `database/sql` with raw SQL (no ORM), `slog` for structured logging
- Mobile-first responsive design
- Lighthouse target: 95+ all categories
- Accessibility: WCAG 2.1 AA minimum
- Three.js effects must be subtle and complement content, not compete with it
- HeroScene adapts opacities for light/dark theme via MutationObserver

## Design System

- **Palette:** Blue primary (#2563EB), teal secondary (#0D9488), with light/dark theme via CSS custom properties and `[data-theme="dark"]`. Blue = action/CTA, Teal = taxonomy/information.
- **Typography:** Satoshi (headings), DM Sans Variable (body), JetBrains Mono Variable (code)
- **Effects:** Glass morphism on cards/quoter, Three.js icosahedron + particles in hero
- **Tone:** Casual-technical, creative, minimalist

## SEO

- JSON-LD structured data (`components/seo/StructuredData.astro`): LocalBusiness, Person, BlogPosting, Service
- Auto-generated sitemap.xml and RSS feed (`/feed.xml`, Spanish posts only)
- Dynamic OG images via satori + sharp (`pages/og/[...slug].png.ts`)
- `llms.txt` at site root for AI discoverability
- Target keywords: "desarrollador web ensenada", "programador baja california", "sistemas administrativos ensenada"

## Docker Notes (Local Development Only)

- `docker-compose.yml` + `docker-compose.dev.yml` are for local development only
- Dev mode disables Traefik/Nginx/Gatus, exposes API ports directly (8081, 8082)
- CORS origin defaults to `http://localhost:4321` in dev
- Alpine runtime images need `tzdata` package for timezone support (America/Tijuana)
- SQLite data persisted via Docker volume `sqlite-data:/data`
- Production uses K3s — see `k8s/` directory

## Content TODOs

- **About section photo:** Replace the placeholder SVG illustration in `components/about/About.astro` with real photos of Joel (speaking at AWS User Group Ensenada, GDG Tijuana, or similar tech events). The current SVG is a generic silhouette and weakens credibility.

## Reference

- Pending tasks: `TODO.md`
