# Frontend — LMS Administrator Panel

React SPA consumed by the library's single Administrator to manage the catalog, students,
loans, returns, and penalties.

## Tech Stack

* **Framework:** React 19 + TypeScript
* **Build tool:** Vite
* **Styling:** Tailwind CSS v4
* **Routing:** React Router
* **HTTP client:** axios

## Structure

```
src/
├── components/
│   ├── layout/     → AppLayout (sidebar/topbar), ProtectedRoute
│   └── ui/         → Button, Card, EmptyState — shared components
├── lib/
│   ├── api.ts       → axios instance (attaches JWT, handles 401)
│   └── auth.ts       → token storage
├── pages/            → one folder per screen area (students/, books/, loans/)
├── types/             → TS types mirroring the OpenAPI contract
└── App.tsx             → route tree (see library-docs/12-ux-ui/navigation-map.md)
```

Placeholder screens show an "Implemented in HU-XX" badge — real content, forms, and data
fetching land on the `feat/HU-XX-...` branch that implements each user story.

## Development

```bash
npm install
cp .env.example .env   # set VITE_API_BASE_URL if not using the default
npm run dev
```

Or, with the rest of the stack, from the repo root:

```bash
docker compose up --build frontend
```

## Correlations

* Navigation map / access rules → `library-docs/12-ux-ui/navigation-map.md`
* Design tokens (colors/typography — finalized here, in code) → `src/index.css`
* API contract this UI calls → `library-docs/07-api/contracts/openapi/library-api.yaml`
