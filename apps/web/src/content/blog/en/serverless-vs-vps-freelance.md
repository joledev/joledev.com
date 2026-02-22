---
title: "Serverless vs. VPS: what nobody tells you when you work alone"
description: "Real costs, cold starts, security, and why I ended up running K3s on a $15 USD VPS instead of Lambda — with actual numbers, diagrams, and a concrete decision framework."
pubDate: 2026-02-15
updatedDate: 2026-02-22
tags: ["serverless", "vps", "architecture", "freelance", "devops", "kubernetes", "go"]
category: "opinion"
lang: "en"
draft: false
---

Last year a fishing cooperative in Ensenada asked me to build a system for logging catches, generating regulatory reports for CONAPESCA, and tracking cold-storage inventory. Thirty users max. My first instinct was to build it all on AWS: Lambda for the API, DynamoDB for data, S3 for documents, API Gateway up front. Modern, scalable, "correct".

Then I did the math. The system would cost more in monthly infrastructure than the cooperative paid for internet. So I deployed it on a $15 VPS and it's been running for months without a single incident.

That experience crystallized something I'd been suspecting: **the infrastructure the industry recommends and the infrastructure a freelancer actually needs are fundamentally different things**. This article is the long version of that argument — with numbers, code, and diagrams of what I actually run.

## The real costs: money and time

The conversation about serverless costs always starts with "Lambda is free up to 1 million requests per month." That's true. It's also irrelevant, because Lambda doesn't run alone.

Let's do the math with a real project: a Go API with 3 endpoints, a relational database, and email sending. Basically what my portfolio runs.

### Serverless (AWS)

| Service | Monthly cost |
|---------|-------------|
| Lambda (50K invocations, 256MB, 200ms avg) | ~$0.20 |
| API Gateway (50K requests) | ~$0.18 |
| RDS PostgreSQL (db.t4g.micro, minimum) | $12.40 |
| NAT Gateway (if Lambda is in VPC) | $32.40 + data |
| CloudWatch Logs (5GB ingestion) | $2.50 |
| Secrets Manager (3 secrets) | $1.20 |
| Route 53 (hosted zone) | $0.50 |
| **Total** | **~$49 - $70** |

Lambda compute is pennies. Everything else is what kills you. The NAT Gateway in particular is absurd: $32 a month for the privilege of letting your function talk to your database inside a VPC. If you don't use a VPC (and make your RDS public), you have a security problem. If you do, you pay the tax.

### VPS (what I actually use)

| Component | Monthly cost |
|-----------|-------------|
| Hetzner CX22 (2 vCPU, 4GB RAM) | $5.39 |
| Domain (.com) | ~$1.00 (prorated) |
| **Total** | **~$6.40** |

On that same VPS I run my portfolio (joledev.com), the quoter API, the scheduler API, a monitoring dashboard with Gatus, and I still have over 3GB of RAM free. I could fit two more client projects before needing an upgrade.

The difference isn't 10x — it's **nearly an order of magnitude**. And that's without counting the cost of my time. Configuring IAM policies, debugging in CloudWatch, managing Lambda concurrency limits... those are hours that simply don't exist on the VPS.

### The cost nobody counts: cold starts

When a Lambda function hasn't been invoked for ~15 minutes, AWS destroys the container. The next invocation has to create a new one, load your code, and initialize the runtime. That's a cold start.

With Go (which compiles to a binary), cold starts hover around **300-500ms**. With Node.js or Python, it goes up to **500ms-1.5s**. With Java or .NET, you can hit **3-5 seconds**.

For a cron that generates reports, it doesn't matter. For an API a user is waiting on, 500ms of overhead is the difference between "fast" and "why is this slow?". There are mitigations — provisioned concurrency, keep-alive functions — but each one adds complexity and cost.

On a VPS, the Go process is already running. It responds in **1-5ms**. No warm-up, no cold start, no variability. Your application's latency is your code's latency, period.

## The architecture: two worlds

Here are both approaches side by side. The first is what I'd have to build on AWS. The second is what actually runs on my VPS.

### Serverless architecture (AWS)

<div class="mermaid">
graph TB
    Client[Client] --> APIGW[API Gateway<br/>~$0.18/50K req]
    APIGW --> WAF[WAF/Throttling]
    WAF --> Lambda1[Lambda: Quoter<br/>Cold start: 300-500ms]
    WAF --> Lambda2[Lambda: Scheduler<br/>Cold start: 300-500ms]
    Lambda1 --> RDS[(RDS PostgreSQL<br/>$12.40/mo minimum)]
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
    IAM[IAM Roles<br/>1 per function + policies] -.-> Lambda1
    IAM -.-> Lambda2
    SM[Secrets Manager<br/>$1.20/mo] -.-> Lambda1
    SM -.-> Lambda2
</div>

Every box is a service to configure, monitor, and pay for. The Lambda function itself is the simplest part of the diagram — everything around it is the real complexity.

### VPS architecture with K3s (what I use)

<div class="mermaid">
graph TB
    Client[Client] --> Traefik[Traefik Ingress<br/>Auto TLS via cert-manager]
    subgraph K3s [K3s — VPS $5.39/mo]
        Traefik --> Web[Pod: Web<br/>Astro + Nginx]
        Traefik --> Quoter[Pod: API Quoter<br/>Go — 15MB RAM]
        Traefik --> Scheduler[Pod: API Scheduler<br/>Go — 20MB RAM]
        Traefik --> Gatus[Pod: Gatus<br/>Monitoring]
        Scheduler --> SQLite[(SQLite WAL<br/>PersistentVolume)]
    end
    GHA[GitHub Actions] -->|build + push| GHCR[GHCR]
    GHCR -.->|pull| K3s
</div>

Everything runs on one machine. Traefik handles TLS and routing. Each service is a pod with its own container. If a pod dies, K3s restarts it. Deploying is push to main → GitHub Actions builds images → pushes to GHCR → SSHs to server → `kubectl rollout restart`. No surprises.

Yes, I said K3s. Kubernetes. After writing "you don't need Kubernetes" in the previous version of this article, I ended up using it. The difference is that K3s doesn't feel like Kubernetes — it's a 50MB binary you install in 30 seconds that gives you container orchestration, health checks, rolling updates, and secret management without the ceremony of an EKS cluster. I use it to run multiple projects on a single server with isolation between them, and because if a service goes down at 3 AM, K3s brings it back up on its own.

## The code: same logic, different ceremony

To make the point concrete, here's the same endpoint implemented for Lambda and for a regular HTTP server with Chi.

### AWS Lambda handler (Go)

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

    // ... identical business logic ...

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

### Chi / net/http handler (Go) — what I actually use

```go
func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
    var body QuoteRequest
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
        return
    }

    // ... same business logic ...

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
```

The business logic is identical. The difference is ceremony: the Lambda handler needs the AWS SDK, manual parsing of the API Gateway event (no `http.Request`), and constructing the response as a struct with separate `StatusCode`, `Body`, and `Headers`. The Chi handler uses Go's standard interface that hasn't changed since 2012.

That difference looks minor in one endpoint. Multiply it by 10 endpoints, add tests (which now need to mock the Lambda context), and the accumulated complexity is real.

### The K8s manifest: this is how simple deploy is

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

That's it. The container, its env vars from a K8s secret, and resource limits. Traefik handles TLS and routing by hostname and path. You write this once, apply it with `kubectl apply -f`, and forget about it.

## Security and scalability

This is the section missing from most comparisons.

### Security

**Serverless has security by default... in theory.** You don't manage the OS, you don't patch servers, you don't configure firewalls. But IAM's security model is so granular that a permission error is almost inevitable. An overly permissive policy (`Action: "*"`, `Resource: "*"`) is a risk. An overly restrictive one and your function can't even read from the database. And IAM error messages are useless — "Access Denied" without telling you which permission is missing.

Secrets in Lambda go in environment variables or Secrets Manager. If you use env vars, anyone with Lambda console access sees them in plain text. If you use Secrets Manager, you pay $0.40/secret/month and add latency to every cold start from the API call.

**On a VPS with K3s**, security is your responsibility — but it's predictable. SSH with key-only auth (no passwords), fail2ban to block brute force, UFW to close unnecessary ports, and secrets live as `kubectl secrets` that never get committed to the repo. The attack surface is one SSH port and ports 80/443 that Traefik exposes. It's simpler to audit because there are fewer moving parts.

Which is harder to secure correctly? Depends on scale. For a freelancer with 3 services, the VPS is simpler and more auditable. For a team of 20 with 50 Lambdas, IAM + AWS security tools makes more sense because it scales with the team.

### Scalability

**Serverless scales automatically.** You get 10,000 concurrent requests, Lambda creates 10,000 containers. Sounds perfect until your relational database can't handle 10,000 simultaneous connections. Connection pooling via RDS Proxy exists but adds another service (and another cost). And there are hard concurrency limits per region — the default is 1,000 concurrent executions. If you exceed them, your requests get rejected with 429.

**On K3s**, scalability is manual but predictable. A Go pod serving 1,000 req/s uses ~50MB of RAM. If you need more, you add replicas (`replicas: 3`) or enable the HorizontalPodAutoscaler. No billing surprises and no concurrency limits blocking you.

The real question: when do you need to scale beyond a single VPS? For a Go API with SQLite, the bottleneck is disk — and even then you can easily hit **5,000-10,000 requests per second** before needing to think about horizontal scaling. None of my freelance clients have ever reached even 5% of that.

## The decision framework

After going back and forth between both worlds, this is the heuristic I use:

**Use serverless when:**
- Traffic is genuinely unpredictable (spikes from 0 to thousands and back to 0)
- The function is isolated and ephemeral (a webhook, a file processor, a cron)
- Your client already has AWS infrastructure and wants to add functionality
- The infrastructure budget is larger than the cost of your time configuring it

**Use a VPS when:**
- Traffic is predictable (internal users, business hours, <1000 req/s)
- You need consistent latency (no cold starts)
- You want full control over the environment and costs
- You work alone or on a small team with no dedicated DevOps
- You need to run multiple projects without the bill multiplying

**The short answer for freelancers:** if your client can't explain to you why they need serverless, they don't. A VPS with containers and good CI/CD covers 90% of the projects you'll encounter.

## What I actually use today

My current stack for a new project: a Hetzner VPS with K3s, Traefik as ingress with auto TLS via cert-manager, microservices in Go with Chi, Astro frontend with Svelte islands, SQLite for persistence, and GitHub Actions that builds images, pushes them to GHCR, and does a rollout restart on the cluster.

The total RAM usage of joledev.com with its 4 pods (web, api-quoter, api-scheduler, gatus) is around 120MB. On a 4GB VPS, that's nothing. I could run 10 projects like this before needing a bigger machine.

Do I use serverless? For a Stripe webhook that processes payments on an e-commerce project I built, yes. For a cron that generates monthly PDF reports, also yes. But those are supplements — not the foundation. The foundation is a server I completely understand, that I can debug with `kubectl logs` at 3 AM, and that costs me less than a coffee per month per project.

The best infrastructure isn't the most modern one. It's the one that lets you sleep at night.
