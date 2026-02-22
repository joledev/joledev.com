---
title: 'Mobile Inventory Management App'
description: 'Cross-platform app for real-time inventory management with barcode scanning and offline sync.'
heroImage: 'https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=800&q=80'
screenshots:
  - 'https://images.unsplash.com/photo-1551650975-87deedd944c3?w=800&q=80'
  - 'https://images.unsplash.com/photo-1555774698-0b77e0d5fac6?w=800&q=80'
tags: ['React Native', 'TypeScript', 'SQLite', 'REST API']
category: 'mobile'
lang: 'en'
featured: false
order: 4
---

Cross-platform mobile app (iOS and Android) built with React Native for warehouse and retail inventory control. Features barcode scanning via device camera, offline-first sync with local SQLite database, and real-time movement reports.

## Key Features

- Barcode and QR code scanning via camera
- Offline mode with automatic sync when connection is restored
- Dashboard with stock levels and minimum alerts
- Movement history (entries, exits, transfers)
- Exportable reports in PDF and Excel
- Multi-warehouse management with user permissions

## Tech Stack

- **Frontend:** React Native with TypeScript
- **Local DB:** SQLite via WatermelonDB
- **Backend:** Node.js REST API with PostgreSQL
- **Auth:** JWT with refresh tokens
- **Notifications:** Firebase Cloud Messaging
