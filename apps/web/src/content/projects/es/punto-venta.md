---
title: 'Sistema Punto de Venta — Restaurantes'
description: 'POS completo para restaurantes con gestión de mesas, cocina en tiempo real, inventario y corte de caja.'
heroImage: 'https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?w=800&q=80'
screenshots:
  - 'https://images.unsplash.com/photo-1556742111-a301076d9d18?w=800&q=80'
tags: ['Electron', 'React', 'SQLite', 'WebSockets', 'Node.js']
category: 'system'
lang: 'es'
featured: false
order: 9
---

Sistema punto de venta de escritorio desarrollado con Electron para cadena de restaurantes. Funciona offline con SQLite local y sincroniza con servidor central cuando hay conexión. Incluye pantalla de cocina en tiempo real, gestión de mesas y corte de caja.

## Funcionalidades principales

- Interfaz táctil optimizada para pantallas de 10" a 15"
- Mapa de mesas interactivo con drag-and-drop
- Pantalla de cocina con tickets en tiempo real (WebSockets)
- División de cuentas y pagos mixtos (efectivo, tarjeta, vales)
- Inventario con alertas de stock mínimo y mermas
- Corte de caja con desglose por método de pago
- Reportes de ventas diarios, semanales y mensuales
- Modo offline con sincronización automática

## Stack técnico

- **Aplicación:** Electron + React con TypeScript
- **Base de datos local:** SQLite con better-sqlite3
- **Servidor central:** Node.js + Express + PostgreSQL
- **Tiempo real:** WebSockets para pantalla de cocina
- **Impresión:** Protocolo ESC/POS para impresoras térmicas
- **Respaldos:** Sincronización incremental con servidor central
