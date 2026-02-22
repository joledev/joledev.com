---
title: "Serverless vs. VPS: what nobody tells you when you work alone"
description: "After years choosing infrastructure for real projects, I share my process for deciding between serverless and a VPS — and why I ended up where I did."
pubDate: 2026-02-15
tags: ["serverless", "vps", "architecture", "freelance", "devops", "go", "svelte"]
category: "opinion"
lang: "en"
draft: false
---

A couple of years ago I sat down to redesign the way I delivered projects. I'd been working with shared hosting, the occasional half-configured VPS, and that constant nagging feeling that "I should be using something more modern." Serverless was everywhere: Lambda, Cloud Functions, managed containers. The narrative was clear — if you're not in the cloud, you're falling behind.

So I did what any curious developer would do: I dove in headfirst. I set up projects on AWS, experimented with Lambda and API Gateway, tried DynamoDB, built pipelines with Step Functions. And I learned a lot. But I also learned something that nobody mentioned in the tutorials or Twitter threads: **that serverless wasn't designed to solve the problems I actually had**.

This article isn't a generic "pros and cons" comparison. It's what I discovered after going back and forth between both worlds, and why my current infrastructure looks the way it does.

## The starting problem

When you work alone or on a very small team, your scarcest resource isn't compute or storage — it's your time. Every hour you spend on infrastructure is an hour you're not building features, serving clients, or simply resting so you don't burn out.

And that's where I started noticing a disconnect between what the industry recommends and what I actually needed. Most articles about cloud architecture are written from the perspective of teams with at least a dedicated DevOps engineer, monthly infrastructure budgets in the four figures, and traffic volumes that justify the complexity. But when your client is a local company with 50 employees that needs an internal system, that context simply doesn't apply.

I'm not saying serverless is bad. I'm saying the right question isn't "serverless or VPS?" but rather "what problem am I solving and what's the most direct way to solve it?"

## What I discovered using serverless on real projects

My first serious project with Lambda was a notification system. The logic was straightforward: receive an event, process data, send an email or a message. On paper, the perfect serverless use case — short functions, event-driven, stateless.

And it worked. But the road to get there was revealing.

Setting up IAM permissions took longer than I expected. Not because it's conceptually difficult, but because AWS's permission model is absurdly granular. Each function needed a role, each role needed specific policies, and an error in any policy resulted in a cryptic message that sent you to Stack Overflow to figure out which permission was missing. For a team with a security engineer who handles that, it's not a problem. For me, alone, it was pure friction.

Then came debugging. When something failed in Lambda, the flow was: check CloudWatch logs (which have their own pricing), find the request ID, locate the error, make a change, deploy, wait for the cold start, and test again. Compare that with having a server where you run `docker logs` and see everything in real time, or set a breakpoint and debug directly.

Cost was also a surprise. Lambda compute is cheap, yes. But API Gateway charges per request. If you need a relational database, RDS has a monthly base cost that doesn't go below $15 USD even if you're not using it. If your Lambda needs to access resources inside a VPC, you need a NAT Gateway — which costs around $30 USD/month just for existing. CloudWatch charges for log ingestion. Suddenly, a project that ran completely on a $10 VPS was costing me $50-70 monthly on AWS. And that's not counting my time configuring everything.

[AWS's pricing documentation](https://aws.amazon.com/lambda/pricing/) is transparent about compute costs, but the collateral costs — the Gateway, the NAT, the logs, artifact storage — are what catch you off guard.

## The moment of reflection

After three or four projects like this, I had to sit down and think honestly. Was I using serverless because it solved a real problem, or because I wanted it on my resume? The answer was uncomfortable.

The reality is that most of my clients have predictable traffic. There are no Black Friday spikes. There are no millions of concurrent requests. There are 30, 50, maybe 200 users accessing an internal system during business hours. For that usage pattern, a server running 24/7 for $15-25 a month isn't waste — it's simplicity.

And simplicity has enormous value that's hard to quantify. When something breaks at 11 PM (because it always breaks at 11 PM), I want to be able to SSH in, check the logs, identify the problem, and fix it. I don't want to open the AWS console, navigate between 6 different services, and discover that the error is in a Step Function step that triggers a Lambda that writes to DynamoDB. I want to see a clear stack trace in one place.

## So, serverless is useless?

That's not what I'm saying. Serverless solves real, concrete problems. It solves them very well. But those problems are specific:

**Workloads that scale from zero to thousands and back to zero.** If you have an event that generates massive traffic for a few hours and then nothing, paying for a 24/7 server would be wasting money. Serverless scales automatically and you only pay for what you use. This is particularly valuable for things like batch file processing, where you receive a large batch of documents once a day and the rest of the time there's nothing to process.

**Isolated functions that don't justify their own server.** A webhook that receives notifications from an external service and processes them. A cron that generates a PDF report every Monday. An endpoint that resizes images on upload. These are functions that run sporadically, last seconds, and don't need to maintain state. This is the sweet spot for Lambda and Cloud Functions.

**Integrations within an existing cloud ecosystem.** If your client already has all their infrastructure on AWS and wants to add new functionality, fighting against the ecosystem makes no sense. You use what's already there.

The problem is when you take these advantages and extrapolate them to your entire architecture. I've seen (and made the mistake of building) entire applications where every endpoint is a Lambda, state is managed between DynamoDB and S3, and communication between services goes through SQS and SNS. It works, but the operational complexity is disproportionate to what the project actually needs.

## The road back to VPS (with lessons learned)

When I decided to re-center my infrastructure on VPS, it wasn't a step backward. It was applying everything I'd learned from serverless — automated deploys, reproducible infrastructure, decoupled services — but on a foundation I could manage alone.

The catalyst was Docker. Everything I liked about the serverless model — packaging my code with its dependencies, deploying reproducibly, scaling services independently — I could do with containers on a VPS. Without the abstraction layer of a cloud provider, without the complex pricing, and with full control over my environment.

My current setup is a VPS with Docker Compose, Traefik as a reverse proxy with automatic SSL via Let's Encrypt, and each service running in its own container. Deploying is a `docker compose up -d --build` after a push to main. GitHub Actions handles CI. There's no magic, no vendor lock-in, and if tomorrow I want to move everything to another VPS provider, it's copying files and running the same command.

Traefik deserves special mention because it solved something that with Nginx required manual configuration: automatic routing based on Docker labels and SSL certificate renewal. Defining that `api.mydomain.com` points to the right service is a label in the `docker-compose.yml`, not a separate config file. [Traefik's documentation](https://doc.traefik.io/traefik/) is some of the best I've seen in infrastructure tooling.

## Choosing the stack: why Go and Svelte

This part of the process was harder than I expected, because it involves killing sacred cows.

I came from PHP and TypeScript. I knew Laravel, I knew Express, I knew Next.js. Why change something that already worked? The reason was purely practical: I wanted to run multiple microservices on a VPS without them eating up all the RAM.

A Laravel service with PHP-FPM uses between 80 and 150 MB of RAM at rest. An Express/Node one between 60 and 120 MB. That's not a problem if you have a monolith. But if your architecture has 4 or 5 independent services, you're already at 400-600 MB just in runtime, not counting the database. On a 2 GB VPS, that's tight.

Go changed that equation. A compiled Go service uses between 10 and 25 MB of RAM serving real traffic. I can run 5 or 6 microservices, a database, a reverse proxy, and a monitoring service on a 2 GB VPS and still have over a gig of RAM left. For a freelancer, that means being able to host several client projects on the same machine if needed, or having enormous room to grow.

But resource efficiency wasn't the only reason. Go has something I didn't find in other backend languages: time-proof readability. When I open a Go file I wrote 8 months ago, I understand what it does in minutes. There are no decorators to investigate, no magic dependency injection, no implicit middleware. The program flow is linear and explicit. You receive a request, validate, process, respond. Error handling with `if err != nil` seems verbose at first, but after a while you realize it's exactly the clarity you want when you're debugging at 11 PM.

The router I use is [Chi](https://github.com/go-chi/chi), which follows Go's standard `net/http` interface. This means I'm not learning a proprietary framework — I'm using the standard library with a bit of convenience on top. If Chi stopped being maintained tomorrow, migrating to another router or even to raw `net/http` would be changing a few lines.

For the frontend, the decision was more personal. I'd been using React for years and the fatigue was real. Not fatigue with the framework itself, but with the ecosystem that changes every six months. Next.js changing its rendering model, the transition to server components, the App Router vs Pages Router debate, hydration issues... too much movement for something that should be stable.

Svelte was a breath of fresh air. A Svelte component is HTML with declarative reactivity. There's no virtual DOM, no hooks with cryptic rules, no `useEffect` with dependency arrays you're never sure are correct. Styles are scoped to the component by default — you write normal CSS and it doesn't leak to other components. And with Svelte 5, the reactivity model with runes (`$state`, `$derived`, `$effect`) is so straightforward that the code almost reads like pseudocode.

Combined with [Astro](https://astro.build) as the site framework, I get the best of both worlds: static pages where I don't need interactivity (blog, landing, informational pages) and Svelte components hydrated only where they're needed (forms, dashboards, interactive elements). The result is a site that loads almost instantly because the HTML comes pre-rendered from the server, and JavaScript only loads where necessary.

## The database: SQLite and when to actually scale to PostgreSQL

Another decision that goes against the current: I use SQLite for most of my projects.

I know it sounds odd. SQLite has a reputation as "the development database" or "the one you use for prototyping." But that reputation is outdated. SQLite in WAL mode (Write-Ahead Logging) handles concurrent reads without issue. For applications with a read-heavy, moderate-write pattern — which is exactly what most internal systems are — it performs as well or better than PostgreSQL, with the advantage that there's no database server to configure, maintain, or that can go down.

The database is a file. You back it up by copying it (or with tools like [Litestream](https://litestream.io) that do continuous replication to S3). You migrate it by copying it. There's no connection pooling to configure, no database users to manage, no ports to expose.

When do I scale to PostgreSQL? When I need heavy concurrent writes from multiple processes, advanced full-text search, specialized data types like indexed JSON or geospatial data, or when data volume exceeds what's reasonable for a single file. For the vast majority of small and mid-sized business projects, that point never comes.

## The deploy: GitHub Actions, Docker, and nothing else

My deploy pipeline is deliberately simple. Push to `main`, GitHub Actions builds the Docker images, pushes them to the registry, and the VPS pulls and restarts the services. No Kubernetes, no Terraform, no Ansible. A `docker-compose.yml` defines the entire infrastructure.

Is this "correct" according to DevOps best practices? Probably not, if we were talking about a company with 20 services and a platform team. But for a freelancer with 3-6 services per project, it's exactly what I need: something I completely understand, that I can debug without documentation, and that doesn't require maintaining auxiliary infrastructure.

A typical Go service Dockerfile is minimal:

```dockerfile
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /server .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

The final image weighs between 15 and 25 MB depending on dependencies. Compare that with a Node image that easily exceeds 200 MB, or a PHP-FPM one that runs around 150 MB. The difference shows in deploy times and how many images you can keep in the registry without storage costs spiraling.

## Conflicts I still have

I don't want to paint this as the perfect solution, because it's not. There are things that still cause friction.

**Go doesn't have an ORM that fully convinces me.** GORM is popular but adds a layer of magic that goes against Go's philosophy. SQLx is more explicit but requires writing SQL by hand for everything. I ended up writing raw SQL with the standard `database/sql` library, and while it works well, I sometimes miss the productivity of Eloquent in Laravel or Prisma in TypeScript. It's a conscious tradeoff: more control and transparency in exchange for more lines of code in the data layer.

**Svelte has a smaller ecosystem than React.** When you need a third-party component — a complex date picker, a rich text editor, a charting library — the options in Svelte are fewer and sometimes less mature. This is becoming less of a problem as the ecosystem grows, but it's a reality to consider.

**Monitoring on a VPS is manual.** On AWS you have CloudWatch, X-Ray, automatic dashboards. On a VPS, if you want monitoring, you set it up yourself. Personally, I use structured logging with `slog` (Go's standard library since version 1.21) and send it to stdout for Docker to capture. It's not as sophisticated as a full APM, but for my needs it's enough. Tools like [Dozzle](https://dozzle.dev/) give you a web-based log dashboard with zero configuration, which helps a lot.

**Horizontal scalability has a ceiling.** If a project grows to the point of needing multiple servers, load balancing, and automatic failover, a single VPS with Docker Compose falls short. But in my experience, that point comes much later than most people think. A well-optimized 4 GB VPS with Go comfortably serves thousands of requests per second. When a project reaches that traffic level, it usually has the budget for more robust infrastructure.

## Current state

Today my standard infrastructure for a new project looks like this: a VPS on Hetzner or DigitalOcean ($15-25 USD/month), Docker Compose for orchestration, Traefik for reverse proxy and SSL, Go backend services with Chi, Astro frontend with Svelte components, SQLite for persistence with Litestream for backups, and GitHub Actions for CI/CD.

The total RAM a typical project uses with 2-3 microservices, the frontend, the database, and Traefik is around 200-300 MB. That leaves me enormous headroom on a 2 GB VPS, and the possibility of hosting more than one project per machine if the clients are small.

Do I use serverless? Yes, but for what makes sense. A webhook that receives data from an external service and processes it. A heavy cron that isn't worth keeping in a running container. Specific functions where the pay-per-invocation model genuinely saves money. But it's not the foundation of my architecture — it's a complement.

## What I'd tell someone who's deciding

If you're starting out as a freelancer or reconsidering your infrastructure, my advice is to resist the pressure to adopt what's trendy and focus on what lets you ship fast and maintain with confidence.

You don't need Kubernetes to serve an application with 200 users. You don't need a distributed microservices architecture on Lambda for a CRUD with reports. You don't need DynamoDB for a table with 10,000 records. But you do need to understand what happens when something breaks, and be able to fix it without consulting three different documentations.

Choose the tool you understand inside and out. If that's PHP with Laravel on shared hosting, it works. If it's Node on a VPS, it works. If it's Go with Docker, it works. The technology matters less than the industry wants you to believe. What matters is that you can deliver, that you can maintain, and that your client is happy.

And if after reading all this you still want to try serverless — do it. But try it on an isolated function in a real project, not on an entire system. That way you learn the model without betting your whole project on an architecture you might not need.

At the end of the day, the best infrastructure is the one that lets you focus on the code and your client's problems, not on the infrastructure itself.
