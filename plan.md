# Onboard — HRMS Plan for SMSE (10–500 employees)

> Target: Small / Micro / Small-Medium Enterprises. Auth/AuthZ excluded — plug-and-play from your own project (this service only validates passthrough headers like `X-Org-ID` / `X-User-ID`).

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

### P3 — Money (first paywall, revenue driver #1)
- Payroll: salary structure, CTC split, per-month run, arrears, payslips (PDF), loans/advances, reimbursements
- Compliance toggle by country:
  - **India:** PF, ESI, PT (state slabs), LWF, TDS, Form 16 / 12BB, Bonus, Gratuity, Shops & Establishments leave rules, challans & returns
  - **Generic:** configurable tax slabs, withholding, audit logs, GDPR-ready exports
- Exports: Tally/QuickBooks, bank transfer file (NEFT/IMPS), UAN upload

### P4 — Growth (second paywall, revenue driver #2)
- Simple ATS (careers page, pipeline kanban, offer → onboarding), no resume-AI in v1
- Basic LMS (courses, certs, completions), performance lite (goals, 360-lite, review cycles), engagement surveys
- Advanced analytics, custom workflows, API/webhooks, BI export

**Cut for now:** succession planning, workforce budgeting, advanced OKR, custom BI, AI-hype, succession — add when a paying SMSE asks.

## 4. Plans & Feature Matrix (per employee/month)

Inspired by Zoho People / greytHR / Keka — payroll is the conversion trigger.

| Capability | Free (≤10) | Starter | Growth |
|---|---|---|---|
| **Best for** | trial / micro | 11–100, single location | 101–500, multi-location/shift |
| **Cap** | 10 emp, 1 admin | 100 emp | 500 emp, then Scale custom |
| **Core HR + Docs + Letters + Org Chart** | ✓ | ✓ | ✓ |
| **Attendance / Leave / Shifts / Holidays** | ✓ basic | ✓ full + roster | ✓ + multi-location, overtime, geo |
| **ESS mobile/web + Directory** | ✓ | ✓ | ✓ |
| **Onboarding / Offboarding / Assets / Helpdesk** | — | ✓ | ✓ |
| **Payroll + Payslips + Reimbursements** | — | ✓ **(core paywall)** | ✓ + advances/loans/arrears |
| **Compliance pack** | — | India base | + state packs, audit reports |
| **Reports / Exports / Bulk CSV** | basic | full | advanced + API |
| **ATS / LMS / Performance lite** | — | — | ✓ **(growth paywall)** |
| **Support** | community | email/WhatsApp | priority + onboarding assist |
| **Pricing cue** | ₹0 | ₹59–99 | ₹149–199 |

Free → Starter conversion = payroll. Starter → Growth = ATS/performance.

## 5. Data Model Sketch (Postgres)

```
orgs → branches → departments → employees (→ employments)
shifts, attendance_punches, leave_types, leave_balances, leave_requests, holidays
salary_structures, payroll_runs, payroll_items, reimbursements
documents, assets, tickets, announcements
jobs, candidates, courses, enrollments, review_cycles (P4)
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
# P3
/api/v1/payroll/runs, /payslips, /reimbursements
# P4
/api/v1/jobs, /candidates, /courses, /reviews
```

Keep `loggingMiddleware` wrapped around `ServeMux`. Tenant via `X-Org-ID` header (no auth logic).

## 7. Roadmap (no hard dates until velocity picked)

- **Phase 0.1 (1w) — Hardening:** `goose` migrations, `sqlc`, pagination+validation middleware, tenant header passthrough, `golangci-lint`, Dockerfile deferred (YAGNI)
- **Phase 1 (4–6w) — MVP:** P0 + P1 + ESS list/views, CSV, manager inbox; demo with 50 fake employees
- **Phase 2 (1–2w) — Ops:** onboarding/checklists + helpdesk + assets
- **Phase 3 (3–4w) — Money:** payroll engine + payslips + bank file (paywall) + India compliance base
- **Phase 4 (4w) — Growth:** ATS-lite + performance-lite + analytics + webhooks
- Always: `internal/<module>/{handler,service,store}` — no interface with one impl, no factory for one product.

## 8. Non-Goals (explicitly excluded)

- Authentication / Authorization / SSO / RBAC enforcement — provided externally; this service trusts `X-User-ID`/`X-Org-ID`.
- Native mobile app (API only; PWA covers ESS v1)
- Resume parsing AI, advanced OKR, custom workflow builder

## 9. Ops & Deployment

- Single Postgres, single binary (`cmd/onboard`), `configs/config.yaml`, `go test ./...`, `go vet`
- Later: Dockerfile, GH Actions, `migrate up`, Prometheus `/metrics` (defer)
- Backup: daily `pg_dump`, point-in-time via WAL

## 10. Monetization Notes

- Flat per-employee/month, annual prepay discount, no per-module fees (SMSE hates nickel-and-dime).
- Free 10 is lead magnet; enforce soft limit (block 11th employee without upgrade).
- Bundle payroll+compliance — SMSE will pay to avoid separate vendor.

## 11. Open Questions (answer when ready)

1. Geography: India-first compliance (PF/ESI/PT) or generic from day one?
2. Payroll: build own engine vs integrate RazorpayX / external API for Phase 3?
3. Mobile: PWA enough for ESS or native wrapper planned?
4. ATS depth: job board + pipeline enough, or need resume parsing in v1?
5. Pricing caps: Free 10 ok, or prefer 5 / 25?

---
*Ponytail note: skipped custom cache/abstraction/config system, RBAC, SSO, LMS editor, resume AI — add when a paying SMSE asks for it.*
