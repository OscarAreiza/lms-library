# Loan Management System (LMS) — `lms-library`

A web-based library management system that centralizes book inventory, student records, loans,
returns, and penalties through a single-administrator interface.

> **Source of truth:** all requirements, domain rules, architecture decisions, and API contracts
> live in the [`library-docs`](https://github.com/code-corhuila/library-docs) repository. This
> repo implements what that one specifies — if code and docs disagree, the docs repo wins until
> this README (or a PR) updates it.

## 🚀 Tech Stack

* **Frontend:** React + Vite + TypeScript + Tailwind CSS
* **Backend:** Go (Golang) — Hexagonal Architecture
* **Database:** PostgreSQL
* **Reverse proxy:** NGINX (production topology — see `library-docs/05-architecture/decisions/records/ADR-003-nginx-reverse-proxy.md`)
* **Infrastructure:** Docker / Docker Compose

## 👥 Team Members

* **Oscar Mauricio Areiza** — Tech Lead
* **Luis Alejandro Meneses** — DevOps
* **Hermes Pascuas Herrera** — DevOps

## 📁 Structure

```
backend/    → Go API (hexagonal: domain / application / infrastructure)
frontend/   → React SPA (Administrator panel)
infra/      → Deployment/reverse-proxy/CI-CD config (added later — see infra/README.md)
```

## 🌿 Branch strategy

```
main   ← Production. Stable, released only.
 └── QA   ← Release candidate, validated before promoting to main.
      └── dev   ← Integration branch. All feature work merges here first.
           └── feat/HU-XX-description
           └── fix/HU-XX-description
           └── chore/description
```

Full convention (naming, commit format, PR policy) →
`library-docs/00-governance/git-conventions.md`. Short version:

- One task = one short-lived branch (`feat/HU-06-loan-registration-form`, not one giant branch
  per user story) = one small PR into `dev`.
- Conventional Commits, referencing the HU: `feat(circulation): register a loan (HU-06)`.
- Every PR uses `.github/PULL_REQUEST_TEMPLATE.md` and links its HU.

## ⚙️ Getting Started (Local Development)

```bash
cp .env.example .env
docker compose up -d --build
curl http://localhost:8080/health
```

See `backend/README.md` and `frontend/README.md` for running each service outside Docker.
