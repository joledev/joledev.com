#!/bin/bash
set -euo pipefail

# Deploy script for joledev.com
# Called by GitHub Actions on push to main
# Can also be run manually on the VPS

PROJECT_DIR="/opt/joledev"
cd "$PROJECT_DIR"

echo "==> Pulling latest code..."
git pull origin main

echo "==> Building Astro frontend (via Docker)..."
docker run --rm \
  -v "$PROJECT_DIR/apps/web:/app" \
  -w /app \
  node:22-alpine \
  sh -c "corepack enable && pnpm install --frozen-lockfile && pnpm build"

echo "==> Building Docker images for APIs..."
docker compose build api-quoter api-scheduler

echo "==> Restarting services..."
docker compose up -d

echo "==> Cleaning up old Docker images..."
docker image prune -f

echo "==> Deploy complete!"
docker compose ps
