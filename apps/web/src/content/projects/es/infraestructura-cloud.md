---
title: 'Infraestructura Cloud AWS — Migración Enterprise'
description: 'Migración de infraestructura on-premise a AWS con IaC, alta disponibilidad y reducción de costos del 40%.'
heroImage: 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=800&q=80'
screenshots:
  - 'https://images.unsplash.com/photo-1544197150-b99a580bb7a8?w=800&q=80'
tags: ['AWS', 'Terraform', 'Docker', 'GitHub Actions', 'Linux']
category: 'devops'
lang: 'es'
featured: false
order: 8
---

Proyecto de migración completa de infraestructura on-premise a AWS para una empresa con más de 200 empleados. Incluyó diseño de arquitectura cloud-native, implementación de Infrastructure as Code con Terraform, y automatización de despliegues con CI/CD.

## Alcance del proyecto

- Migración de 12 servidores físicos a servicios AWS
- Diseño de VPC multi-AZ con alta disponibilidad
- Implementación de auto-scaling para cargas variables
- Pipeline CI/CD completo (build, test, deploy, rollback)
- Monitoreo y alertas con CloudWatch y PagerDuty
- Reducción de costos operativos del 40%

## Stack técnico

- **Cloud:** AWS (EC2, ECS, RDS, S3, CloudFront, Lambda, SQS)
- **IaC:** Terraform con módulos reutilizables
- **Contenedores:** Docker + ECS Fargate
- **CI/CD:** GitHub Actions con ambientes staging/production
- **Monitoreo:** CloudWatch, Grafana Cloud, PagerDuty
- **Seguridad:** AWS WAF, Security Groups, IAM policies, Secrets Manager
