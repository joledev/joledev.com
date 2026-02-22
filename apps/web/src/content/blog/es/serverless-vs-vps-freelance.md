---
title: "Serverless vs. VPS: lo que nadie te dice cuando trabajas solo"
description: "Costos reales, cold starts, seguridad, y por qué terminé corriendo K3s en un VPS de $15 USD en vez de usar Lambda — con números, diagramas y un framework de decisión concreto."
pubDate: 2026-02-15
updatedDate: 2026-02-22
tags: ["serverless", "vps", "arquitectura", "freelance", "devops", "kubernetes", "go"]
category: "opinion"
lang: "es"
draft: false
---

El año pasado una cooperativa pesquera en Ensenada me pidió un sistema para registrar capturas, generar reportes para CONAPESCA, y controlar inventario de producto en cámara fría. Treinta usuarios máximo. Mi primer instinto fue armar todo en AWS: Lambda para la API, DynamoDB para los datos, S3 para los documentos, API Gateway al frente. Moderno, escalable, "correcto".

Después hice las cuentas. El sistema iba a costar más en infraestructura mensual que lo que la cooperativa pagaba de internet. Así que lo deployé en un VPS de $15 dólares y lleva meses corriendo sin un solo incidente.

Esa experiencia cristalizó algo que venía sospechando: **la infraestructura que la industria recomienda y la que un freelance necesita son cosas fundamentalmente diferentes**. Este artículo es la versión larga de ese argumento — con números, código, y diagramas de lo que realmente uso.

## Los costos reales: dinero y tiempo

La conversación sobre costos de serverless siempre empieza con "Lambda es gratis hasta 1 millón de requests al mes". Eso es cierto. También es irrelevante, porque Lambda no corre sola.

Vamos a hacer las cuentas con un proyecto real: una API en Go con 3 endpoints, base de datos relacional, y envío de emails. Lo que tiene mi portafolio, básicamente.

### Serverless (AWS)

| Servicio | Costo mensual |
|----------|--------------|
| Lambda (50K invocaciones, 256MB, 200ms promedio) | ~$0.20 |
| API Gateway (50K requests) | ~$0.18 |
| RDS PostgreSQL (db.t4g.micro, mínimo) | $12.40 |
| NAT Gateway (si Lambda está en VPC) | $32.40 + datos |
| CloudWatch Logs (5GB ingesta) | $2.50 |
| Secrets Manager (3 secrets) | $1.20 |
| Route 53 (zona hospedada) | $0.50 |
| **Total** | **~$49 - $70** |

El compute de Lambda es centavos. Todo lo demás es lo que te mata. El NAT Gateway en particular es absurdo: $32 dólares al mes por el privilegio de que tu función pueda hablar con tu base de datos dentro de una VPC. Si no usas VPC (y pones tu RDS público), tienes un problema de seguridad. Si sí la usas, pagas el impuesto.

### VPS (lo que realmente uso)

| Componente | Costo mensual |
|------------|--------------|
| Hetzner CX22 (2 vCPU, 4GB RAM) | $5.39 |
| Dominio (.com) | ~$1.00 (prorrateado) |
| **Total** | **~$6.40** |

En ese mismo VPS corro mi portafolio (joledev.com), la API del cotizador, la API del agendador, un dashboard de monitoreo con Gatus, y todavía me sobra más de 3GB de RAM. Podría meter dos proyectos más de clientes antes de necesitar un upgrade.

La diferencia no es 10x — es **casi un orden de magnitud**. Y eso sin contar el costo de mi tiempo. Configurar IAM policies, debuggear en CloudWatch, manejar los límites de concurrencia de Lambda... son horas que en el VPS simplemente no existen.

### El costo que nadie cuenta: cold starts

Cuando una función Lambda no ha sido invocada en ~15 minutos, AWS destruye el contenedor. La siguiente invocación tiene que crear uno nuevo, cargar tu código, e inicializar el runtime. Eso es el cold start.

Con Go (que compila a binario), los cold starts rondan los **300-500ms**. Con Node.js o Python, sube a **500ms-1.5s**. Con Java o .NET, puedes llegar a **3-5 segundos**.

Para un cron que genera reportes, no importa. Para una API que un usuario está esperando, 500ms de overhead es la diferencia entre "rápido" y "¿por qué tarda?". Hay mitigaciones — provisioned concurrency, funciones keep-alive — pero cada una agrega complejidad y costo.

En un VPS, el proceso de Go ya está corriendo. Responde en **1-5ms**. No hay warm-up, no hay cold start, no hay variabilidad. La latencia de tu aplicación es la latencia de tu código, punto.

## La arquitectura: dos mundos

Estos son los dos enfoques lado a lado. El primero es lo que tendría que armar en AWS. El segundo es lo que realmente corre en mi VPS.

### Arquitectura serverless (AWS)

<div class="mermaid">
graph TB
    Client[Cliente] --> APIGW[API Gateway<br/>~$0.18/50K req]
    APIGW --> WAF[WAF/Throttling]
    WAF --> Lambda1[Lambda: Quoter<br/>Cold start: 300-500ms]
    WAF --> Lambda2[Lambda: Scheduler<br/>Cold start: 300-500ms]
    Lambda1 --> RDS[(RDS PostgreSQL<br/>$12.40/mo mínimo)]
    Lambda2 --> RDS
    Lambda1 --> SES[SES - Email]
    Lambda2 --> SES
    Lambda1 --> CW[CloudWatch<br/>$2.50/mo logs]
    Lambda2 --> CW
    subgraph VPC [VPC - NAT Gateway $32.40/mo]
        Lambda1
        Lambda2
        RDS
    end
    IAM[IAM Roles<br/>1 por función + políticas] -.-> Lambda1
    IAM -.-> Lambda2
    SM[Secrets Manager<br/>$1.20/mo] -.-> Lambda1
    SM -.-> Lambda2
</div>

Cada caja es un servicio que configurar, monitorear, y pagar. La función Lambda en sí es lo más simple del diagrama — todo lo que la rodea es la complejidad real.

### Arquitectura VPS con K3s (lo que uso)

<div class="mermaid">
graph TB
    Client[Cliente] --> Traefik[Traefik Ingress<br/>TLS automático vía cert-manager]
    subgraph K3s [K3s — VPS $5.39/mo]
        Traefik --> Web[Pod: Web<br/>Astro + Nginx]
        Traefik --> Quoter[Pod: API Quoter<br/>Go — 15MB RAM]
        Traefik --> Scheduler[Pod: API Scheduler<br/>Go — 20MB RAM]
        Traefik --> Gatus[Pod: Gatus<br/>Monitoreo]
        Scheduler --> SQLite[(SQLite WAL<br/>PersistentVolume)]
    end
    GHA[GitHub Actions] -->|build + push| GHCR[GHCR]
    GHCR -.->|pull| K3s
</div>

Todo corre en una máquina. Traefik se encarga del TLS y el ruteo. Cada servicio es un pod con su propio contenedor. Si un pod muere, K3s lo reinicia. El deploy es push a main → GitHub Actions construye las imágenes → las sube a GHCR → SSH al servidor → `kubectl rollout restart`. Sin sorpresas.

Sí, dije K3s. Kubernetes. Después de escribir "no necesitas Kubernetes" en la versión anterior de este artículo, terminé usándolo. La diferencia es que K3s no se siente como Kubernetes — es un binario de 50MB que instalas en 30 segundos y que te da orquestación de contenedores, health checks, rolling updates, y manejo de secrets sin la ceremonia de un cluster EKS. Lo uso para poder correr múltiples proyectos en un solo servidor con aislamiento entre ellos, y porque si un servicio se cae a las 3 AM, K3s lo levanta solo.

## El código: misma lógica, diferente ceremonia

Para hacer el punto concreto, aquí está el mismo endpoint implementado para Lambda y para un servidor HTTP normal con Chi.

### Handler para AWS Lambda (Go)

```go
package main

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

type QuoteRequest struct {
    ProjectTypes []string `json:"projectTypes"`
    Contact      struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    } `json:"contact"`
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var body QuoteRequest
    if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusBadRequest,
            Body:       `{"error":"invalid JSON"}`,
            Headers:    map[string]string{"Content-Type": "application/json"},
        }, nil
    }

    // ... lógica de negocio idéntica ...

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       `{"success":true}`,
        Headers:    map[string]string{"Content-Type": "application/json"},
    }, nil
}

func main() {
    lambda.Start(handler)
}
```

### Handler para Chi / net/http (Go) — lo que realmente uso

```go
func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
    var body QuoteRequest
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
        return
    }

    // ... misma lógica de negocio ...

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
```

La lógica de negocio es idéntica. La diferencia es la ceremonia: el handler de Lambda necesita el SDK de AWS, parsear el evento de API Gateway manualmente (no hay `http.Request`), construir el response como un struct con `StatusCode`, `Body`, y `Headers` por separado. El handler de Chi usa la interfaz estándar de Go que es la misma desde 2012.

Esa diferencia parece menor en un endpoint. Multiplícala por 10 endpoints, agrega los tests (que ahora necesitan mockear el contexto de Lambda), y la complejidad acumulada es real.

### El manifiesto de K8s: así de simple es el deploy

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-quoter
  namespace: joledev
spec:
  replicas: 1
  selector:
    matchLabels:
      app: api-quoter
  template:
    spec:
      containers:
        - name: api-quoter
          image: ghcr.io/joledev/joledev-api-quoter:latest
          ports:
            - containerPort: 8081
          env:
            - name: RESEND_API_KEY
              valueFrom:
                secretKeyRef:
                  name: joledev-secrets
                  key: RESEND_API_KEY
          resources:
            requests:
              cpu: 10m
              memory: 32Mi
            limits:
              cpu: 200m
              memory: 128Mi
```

Eso es todo. El contenedor, sus variables de entorno desde un secret de K8s, y los límites de recursos. Traefik se encarga del TLS y el ruteo por hostname y path. Escribes esto una vez, lo aplicas con `kubectl apply -f`, y te olvidas.

## Seguridad y escalabilidad

Esta es la sección que falta en la mayoría de las comparaciones.

### Seguridad

**Serverless tiene seguridad por defecto... en teoría.** No manejas el OS, no parcheas servidores, no configuras firewalls. Pero el modelo de seguridad de IAM es tan granular que un error de permisos es casi inevitable. Una policy demasiado permisiva (`Action: "*"`, `Resource: "*"`) es un riesgo. Una demasiado restrictiva y tu función no puede ni leer de la base de datos. Y los mensajes de error de IAM son inútiles — "Access Denied" sin decirte qué permiso falta.

Los secrets en Lambda van en variables de entorno o en Secrets Manager. Si usas env vars, cualquiera con acceso a la consola de Lambda los ve en texto plano. Si usas Secrets Manager, pagas $0.40/secret/mes y agregas latencia a cada cold start por la llamada al API.

**En un VPS con K3s**, la seguridad es tu responsabilidad — pero es predecible. SSH con llave (no contraseña), fail2ban para bloquear brute force, UFW para cerrar puertos innecesarios, y los secrets viven como `kubectl secrets` que nunca se commitean al repo. El surface de ataque es un puerto SSH y los puertos 80/443 que Traefik expone. Es más simple de auditar porque hay menos partes.

¿Cuál es más difícil de asegurar correctamente? Depende de la escala. Para un freelance con 3 servicios, el VPS es más simple y más auditabe. Para un equipo de 20 con 50 Lambdas, IAM + AWS security tools tiene más sentido porque escala con el equipo.

### Escalabilidad

**Serverless escala automáticamente.** Recibes 10,000 requests concurrentes, Lambda crea 10,000 contenedores. Suena perfecto hasta que tu base de datos relacional no puede manejar 10,000 conexiones simultáneas. El connection pooling con RDS Proxy existe pero agrega otro servicio (y otro costo). Y hay límites duros de concurrencia por región — el default es 1,000 ejecuciones concurrentes. Si los excedes, tus requests se rechazan con 429.

**En K3s**, la escalabilidad es manual pero predecible. Un pod de Go sirviendo 1,000 req/s usa ~50MB de RAM. Si necesitas más, agregas réplicas (`replicas: 3`) o activas el HorizontalPodAutoscaler. No hay sorpresas en la factura ni límites de concurrencia que te bloqueen.

La pregunta real: ¿cuándo necesitas escalar más allá de un VPS? Para una API de Go con SQLite, el bottleneck es el disco — y aún así llegas fácilmente a **5,000-10,000 requests por segundo** antes de necesitar pensar en escalar horizontalmente. Ninguno de mis clientes freelance ha llegado ni al 5% de eso.

## El framework de decisión

Después de ir y venir entre ambos mundos, esta es la heurística que uso:

**Usa serverless cuando:**
- El tráfico es genuinamente impredecible (picos de 0 a miles y vuelta a 0)
- La función es aislada y efímera (un webhook, un procesador de archivos, un cron)
- Tu cliente ya tiene infraestructura en AWS y quiere agregar funcionalidad
- El presupuesto de infraestructura es mayor al costo de tu tiempo configurándola

**Usa un VPS cuando:**
- El tráfico es predecible (usuarios internos, horario de oficina, <1000 req/s)
- Necesitas latencia consistente (sin cold starts)
- Quieres control total sobre el entorno y los costos
- Trabajas solo o en equipo chico y no tienes un DevOps dedicado
- Necesitas correr múltiples proyectos sin que la factura se multiplique

**La respuesta corta para freelancers:** si tu cliente no te puede explicar por qué necesita serverless, no lo necesita. Un VPS con contenedores y un buen CI/CD cubre el 90% de los proyectos que vas a encontrar.

## Lo que realmente uso hoy

Mi stack actual para un proyecto nuevo: un VPS en Hetzner con K3s, Traefik como ingress con TLS automático vía cert-manager, microservicios en Go con Chi, frontend en Astro con islas de Svelte, SQLite para persistencia, y GitHub Actions que construye las imágenes, las sube a GHCR, y hace rollout restart en el cluster.

El consumo total de RAM de joledev.com con sus 4 pods (web, api-quoter, api-scheduler, gatus) ronda los 120MB. En un VPS de 4GB, eso es nada. Podría correr 10 proyectos como este antes de necesitar una máquina más grande.

¿Uso serverless? Para un webhook de Stripe que procesa pagos de un e-commerce que hice, sí. Para un cron que genera PDFs de reportes mensuales, también. Pero son complementos — no la base. La base es un servidor que entiendo completamente, que puedo debuggear con `kubectl logs` a las 3 AM, y que me cuesta menos que un café al mes por proyecto.

La mejor infraestructura no es la más moderna. Es la que te deja dormir tranquilo.
