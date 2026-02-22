---
title: 'App Móvil de Control de Inventario'
description: 'Aplicación multiplataforma para gestionar inventario en tiempo real con escaneo de códigos de barras y sincronización offline.'
heroImage: 'https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=800&q=80'
screenshots:
  - 'https://images.unsplash.com/photo-1551650975-87deedd944c3?w=800&q=80'
  - 'https://images.unsplash.com/photo-1555774698-0b77e0d5fac6?w=800&q=80'
tags: ['React Native', 'TypeScript', 'SQLite', 'REST API']
category: 'mobile'
lang: 'es'
featured: false
order: 4
---

Aplicación móvil multiplataforma (iOS y Android) desarrollada con React Native para el control de inventario en almacenes y tiendas. Integra escaneo de códigos de barras mediante la cámara del dispositivo, sincronización offline-first con base de datos local SQLite, y reportes de movimientos en tiempo real.

## Funcionalidades principales

- Escaneo de códigos de barras y QR con la cámara
- Modo offline con sincronización automática al recuperar conexión
- Dashboard con niveles de stock y alertas de mínimos
- Historial de movimientos (entradas, salidas, transferencias)
- Generación de reportes exportables en PDF y Excel
- Gestión multialmacén con permisos por usuario

## Stack técnico

- **Frontend:** React Native con TypeScript
- **Base de datos local:** SQLite vía WatermelonDB
- **Backend:** API REST en Node.js con PostgreSQL
- **Autenticación:** JWT con refresh tokens
- **Notificaciones:** Firebase Cloud Messaging
