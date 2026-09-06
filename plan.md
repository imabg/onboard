# Onboard — HRMS Plan for SMSE (10–500 employees)

> Target: Global, India-first. Small / Micro / Small-Medium Enterprises. Auth/AuthZ excluded — plug-and-play from your own project (this service only validates passthrough headers like `X-Org-ID` / `X-User-ID`).

## 1. Target & Principles

- **Persona:** 1–2 HR generalists, no IT team, 10–500 headcount, WhatsApp-first, needs DIY setup in <1 day.
- **Pain:** spreadsheets for attendance/leaves/payroll, state-wise compliance confusion, no self-service.
- **Principles:** boring over clever, Excel in/out always, mobile ESS first, payroll+attendance+leave bundled, opinionated defaults > configurability. Ship the MVP 80% can use; hide power features behind Growth.

## 2. Current Foundation (already in repo)

- Go 1.27 + `net/http ServeMux` (stdlib, method patterns `GET /health`) + `pgx/pgxpool` + `zap` + `viper` + Postgres
- `GET /health`, `GET /ready`, `GET /api/v1/` + `loggingMiddleware`
- Extend same stack — no new deps until proven. `ponytail: global pgx pool, per-tenant pools if noisy-neighbor matters`

## 3. Modules (SMSE-cut)

### P0 — Core HR (MVP, must-have)
- Employee master, org / dept / designation / branch, lifecycle states (probation/active/exit)
- Document vault (Aadhaar/PAN/UAn, offer/relieving), letter generator (offer, relieving, experience)
- Org chart, directory, global search, audit log, soft-delete + history
- Bulk CSV import/export (always), Excel template download

### P1 — Time (MVP, must-have)
- Attendance: web/mobile punch, geo-fence, IP fence, regularization, overtime, bulk biometric import (CSV)
- Shifts & roster, weekly-off, holidays (state-wise), leave types / balances / carry-forward / encashment, comp-off
- 1–2 level approvals (manager → HR), Slack/WhatsApp webhook (defer)

### P2 — ESS & Ops (MVP)
- Employee self-service (profile, payslips, leave/attendance, helpdesk tickets), manager inbox, announcements
- Onboarding/offboarding checklists, asset assignment (laptop/ID), e-sign deferred to integration
- File manager, search, notifications

### P3 — Money (first paywall, revenue driver #1) — RazorpayX integrated
- Payroll via **RazorpayX Payroll / Payouts** (source of truth for salary, compliance, disbursement). Onboard owns: salary structure, CTC split, attendance→payroll inputs, mapping `employee → RazorpayX contact/fund account`, idempotent `payroll_runs` sync.
- Keep a thin local `payroll_runs / payroll_items / payslips` mirror for reports/search; RazorpayX is disbursement + PF/ESI/PT/TDS calculation. RazorpayX webhooks → `payroll_runs.status` (pending→processing→paid/failed).
- Compliance: **India-first** (PF, ESI, PT state slabs, LWF, TDS, Form 16/12BB, Bonus, Gratuity, S&E leave, challans), with global toggle for later (configurable tax slabs, withholding, GDPR exports).
- Exports: Tally/QuickBooks, RazorpayX payout reports, UAN upload. `ponytail: no own payroll engine — RazorpayX owns compliance math, swap only if scale demands`

### P4 — Growth (second paywall, revenue driver #2) — ATS Depth
- **Full ATS (depth):** careers page, requisitions/approvals, pipeline kanban, **resume parsing** (pdf/docx → structured), deduplication, search (skills/exp/CTC), interview scheduling + scorecards, offer letter + e-sign → onboarding handoff. Keep AI optional/pluggable.
- Basic LMS (courses, certs, completions), performance lite (goals, 360-lite, review cycles), engagement surveys
- Advanced analytics, custom workflows, API/webhooks, BI export

**Cut for now:** succession planning, workforce budgeting, advanced OKR, custom BI, AI-hype, succession — add when a paying SMSE asks.

## 4. Plans & Feature Matrix (per employee/month)

Global pricing, India-first compliance. Inspired by Zoho People / greytHR / Keka — payroll (via RazorpayX) is the conversion trigger.

| Capability | Free (≤5) | Starter | Growth |
|---|---|---|---|
| **Best for** | trial / micro | 6–100, single location | 101–500, multi-location/shift |
| **Cap** | 5 emp, 1 admin | 100 emp | 500 emp, then Scale custom |
| **Core HR + Docs + Letters + Org Chart** | ✓ | ✓ | ✓ |
| **Attendance / Leave / Shifts / Holidays** | ✓ basic | ✓ full + roster | ✓ + multi-location, overtime, geo |
| **ESS (Flutter hybrid + web) + Directory** | ✓ | ✓ | ✓ |
| **Onboarding / Offboarding / Assets / Helpdesk** | — | ✓ | ✓ |
| **Payroll (RazorpayX) + Payslips + Reimbursements** | — | ✓ **(core paywall)** | ✓ + advances/loans/arrears |
| **Compliance pack** | — | India base | + state packs, audit reports; global toggle later |
| **Reports / Exports / Bulk CSV** | basic | full | advanced + API |
| **ATS Depth / LMS / Performance lite** | — | — | ✓ **(growth paywall)** |
| **Support** | community | email/WhatsApp | priority + onboarding assist |
| **Pricing cue** | ₹0 | ₹59–99 | ₹149–199 |

Free (5) → Starter conversion = payroll. Starter → Growth = ATS depth/performance.

## 5. Data Model Sketch (Postgres)

```
orgs → branches → departments → employees (→ employments)
shifts, attendance_punches, leave_types, leave_balances, leave_requests, holidays
salary_structures, payroll_runs, payroll_items, reimbursements, razorpayx_payouts + webhook_events
documents, assets, tickets, announcements
jobs, candidates, candidate_sources, resumes_parsed, interviews, scorecards, courses, enrollments, review_cycles (P4)
```

- All business tables have `org_id` tenant column, `created_at/updated_at`, soft-delete.
- Migrations via `goose`, queries via `sqlc` or raw `pgx` (keep deps low).
- Pagination: `?page&per_page` + `Link` header, validation middleware.

## 6. API Shape (versioned)

```
/health, /ready, /api/v1/ (existing)
# P0
/api/v1/employees, /departments, /documents
# P1
/api/v1/attendance/punches, /shifts, /leaves, /holidays
# P2
/api/v1/onboarding/checklists, /assets, /tickets
# P3 (RazorpayX)
/api/v1/payroll/runs, /payslips, /reimbursements, /payroll/webhooks/razorpayx
# P4 (ATS depth)
/api/v1/jobs, /candidates, /candidates/parse, /interviews, /offers, /courses, /reviews
# Mobile (Flutter hybrid)
# Single API — Flutter hits same /api/v1/* with X-Org-ID; push via FCM
```

Keep `loggingMiddleware` wrapped around `ServeMux`. Tenant via `X-Org-ID` header (no auth logic).

## 7. Roadmap (no hard dates until velocity picked)

- **Phase 0.1 (1w) — Hardening:** `goose` migrations, `sqlc`, pagination+validation middleware, tenant header passthrough, `golangci-lint`, Dockerfile deferred (YAGNI)
- **Phase 1 (4–6w) — MVP:** P0 + P1 + ESS list/views (web + Flutter shell), CSV, manager inbox; demo with 50 fake employees
- **Phase 2 (1–2w) — Ops:** onboarding/checklists + helpdesk + assets
- **Phase 3 (3–4w) — Money:** RazorpayX integration (contacts/payouts/webhooks) + payslips + India compliance base (paywall); no own engine
- **Phase 4 (4–6w) — Growth:** ATS-depth (parsing, dedup, scheduling, scorecards) + LMS + performance-lite + analytics + webhooks; ATS is longer than lite — budget 6w
- **Flutter hybrid:** single Flutter codebase (iOS/Android/Web PWA fallback) hitting same API; FCM push, offline punch queue. Keep Go backend as source of truth.
- Always: `internal/<module>/{handler,service,store}` — no interface with one impl, no factory for one product.

## 8. Non-Goals (explicitly excluded)

- Authentication / Authorization / SSO / RBAC enforcement — provided externally; this service trusts `X-User-ID`/`X-Org-ID`.
- Own payroll compliance engine (RazorpayX owns it)
- Advanced OKR, custom workflow builder (defer)

## 9. Ops & Deployment

- Single Postgres, single binary (`cmd/onboard`), `configs/config.yaml`, `go test ./...`, `go vet`
- RazorpayX: store `razorpayx_contact_id`, `fund_account_id` per employee; idempotency keys per run; webhook HMAC verify; DLQ for failed payouts
- Flutter: one repo `mobile/` (or separate) — same API, FCM, geo/attendance uses OS location; offline queue syncs on reconnect
- Later: Dockerfile, GH Actions, `migrate up`, Prometheus `/metrics` (defer)
- Backup: daily `pg_dump`, point-in-time via WAL

## 10. Monetization Notes

- Flat per-employee/month, annual prepay discount, no per-module fees (SMSE hates nickel-and-dime).
- Free **5** is lead magnet; enforce soft limit (block 6th employee without upgrade) — lower support cost vs 10.
- Bundle payroll (RazorpayX) + India compliance — SMSE will pay to avoid separate vendor.
- Global later: price in local currency, keep RazorpayX for India, pluggable provider interface for other regions (defer).

## 11. Decisions Locked (from you)

1. Geography: **Global, India-first** — build state-wise leave/PT/LWF from day one, global toggle as config.
2. Payroll: **RazorpayX integrated** — no own engine.
3. Mobile: **Flutter hybrid** (iOS/Android, single codebase).
4. ATS: **Depth** — parsing, dedup, search, scheduling, scorecards in Growth.
5. Pricing: **Free till 5** (not 10).

---
*Ponytail note: skipped own payroll engine, RBAC, custom cache/abstraction — RazorpayX + Flutter hybrid are now explicit choices.*
