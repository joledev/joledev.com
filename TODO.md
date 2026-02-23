# TODO — joledev.com

Roadmap and pending improvements.

## Scheduler

- [ ] Add admin page in English (`/en/admin/schedule/`)
- [ ] Full end-to-end booking flow testing

## Blog

- [ ] Write more articles (currently only 1: "Serverless vs. VPS")
- [ ] Verify OG images generate correctly for published posts

## Security

- [ ] Set real Cloudflare Turnstile keys in production (`PUBLIC_TURNSTILE_SITE_KEY` in Astro build env, `TURNSTILE_SECRET_KEY` in K8s `joledev-secrets` and rebuild/restart both APIs)

## Frontend / UX

- [ ] Lighthouse audit: target 95+ in all categories
- [ ] Accessibility audit: WCAG 2.1 AA — contrast, aria labels, keyboard navigation
- [ ] Test responsive on real devices (mobile, tablet)

## SEO

- [ ] Validate JSON-LD on individual blog pages (BlogPosting schema)
- [ ] Validate meta tags and OG images with debugging tools

## Content

- [ ] Review 12 project case studies for accuracy
- [ ] Review pricing and multipliers in `lib/quoter-config.ts`
