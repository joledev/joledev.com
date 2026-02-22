---
title: 'Point of Sale System — Restaurants'
description: 'Complete POS for restaurants with table management, real-time kitchen display, inventory, and cash register.'
heroImage: 'https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?w=800&q=80'
screenshots:
  - 'https://images.unsplash.com/photo-1556742111-a301076d9d18?w=800&q=80'
tags: ['Electron', 'React', 'SQLite', 'WebSockets', 'Node.js']
category: 'system'
lang: 'en'
featured: false
order: 9
---

Desktop point-of-sale system built with Electron for a restaurant chain. Works offline with local SQLite and syncs with central server when connected. Includes real-time kitchen display, table management, and cash register closing.

## Key Features

- Touch-optimized interface for 10" to 15" screens
- Interactive table map with drag-and-drop
- Real-time kitchen display with ticket system (WebSockets)
- Split checks and mixed payments (cash, card, vouchers)
- Inventory with minimum stock alerts and waste tracking
- Cash register closing with payment method breakdown
- Daily, weekly, and monthly sales reports
- Offline mode with automatic sync

## Tech Stack

- **Application:** Electron + React with TypeScript
- **Local DB:** SQLite with better-sqlite3
- **Central Server:** Node.js + Express + PostgreSQL
- **Real-time:** WebSockets for kitchen display
- **Printing:** ESC/POS protocol for thermal printers
- **Backups:** Incremental sync with central server
