# UI

React + TypeScript frontend for the pipeline app. Vite, Tailwind, react-hook-form + zod for the create-job form, react-router for the two pages (job list and job detail).

## Running it

```bash
npm install
npm run dev
```

By default it talks to the backend at `http://localhost:8081/api/v1` with the API key `secret-pipeline-key`. Heads up: that `:8081` is just this app's own hardcoded fallback and doesn't actually match the backend's own default of `:8080` — either run the backend with `PORT=8081`, or set `VITE_API_BASE_URL` (and `VITE_API_KEY` if you've changed `API_KEY` on the backend) to whatever your backend is actually on, e.g. in a `.env.local`.

## Scripts

- `npm run dev` — Vite dev server
- `npm run build` — production build (also what the Docker image uses)
- `npm test` — vitest
- `npm run lint` — eslint

## Layout

- `pages/` — `JobList` and `JobDetail`, the only two routes. Both are lazy-loaded from `App.tsx`.
- `components/common/` — shared building blocks: `AppButton`, `AppInput`, `ConfirmModal`, `StatusBadge`, `ErrorBoundary`.
- `components/jobList/` — the create-job modal and its form sections (sources/transforms/aggregations/exports). Only ever rendered from the job list page, and lazy-loaded too since react-hook-form/zod are a decent chunk of JS most visitors landing on the list don't need.
- `components/jobDetail/` — job detail page's sub-views: metric cards, results table, errors table.
- `service/jobService.ts` — the only place in the app that talks to the backend.
- `schemas/jobSpec.ts` — zod schema backing the create-job form's validation.
- `constants/` — all user-facing copy lives here rather than inline in JSX.
- `types/` — mirrors the Go structs on the backend; keep these in sync if `internal/models`/`internal/pipeline` change field names.

## Notes

- Jobs don't have a real "name" field on the backend. The list and detail pages fall back to showing the first export's file path as a stand-in title (`utils/jobTitle.ts`), or the creation date if a job somehow has no exports.
- The axios instance (`config/axios.ts`) has a response interceptor that toasts every failed request automatically — you generally don't need to handle error UI yourself in a component, just let the request reject.
