## Frontend

React + TypeScript + Vite single-page application styled with Tailwind CSS.

### Stack

- **React 19 + TypeScript**
- **Vite** for dev server and build
- **Tailwind CSS** for styling
- **React Router** for routing
- **React Query** for data fetching and caching
- **react-i18next** with JSON locale files

### Project structure (high level)

- `src/App.tsx` – main app shell and routing entry.
- `src/pages/*` – top-level pages:
  - `Login`, `Register`, `Dashboard`, `Settings`, `ChangePassword`, `ResetPassword`, `Confirm`, `NotFound`.
- `src/components/*` – reusable UI components (cards, forms, language selector, logo, messages, theme toggle, etc.).
- `src/hooks/*` – feature-specific hooks:
  - `useLogin`, `useRegister`, `useConfirm`, `useResetPassword`, `useChangePassword`, `useSettings`.
- `src/providers/*` – context providers:
  - `AuthProvider`, `ConfigProvider`, `LangProvider`, `LoaderProvider`, `MessageProvider`, `AppThemeProvider`, `AppProviders`.
- `src/api/*` – API client, types, and codes shared across hooks/pages.
- `src/locales/*` – translation files (`en`, `de`, `pl`, `ua`).

### Running in development

```bash
cd frontend
npm install
npm run dev
```

- Default dev server port is `5174` (can be overridden via CLI).
- In the full stack, Nginx (from Docker) proxies:
  - `/api/` → Go backend
  - `/` → Vite dev server (`http://host.docker.internal:5174` in dev compose)

Make sure the frontend API base URL matches how you run the backend (either via Nginx proxy or directly).

### Build and preview

```bash
cd frontend
npm run build
npm run preview
```

- Production builds are used by the `headless-frontend` image to publish static files that are then served behind Nginx in prod.

### Testing, linting, formatting

- **Tests** (Vitest + Testing Library):

  ```bash
  cd frontend
  npm test
  ```

- **Lint**:

  ```bash
  cd frontend
  npm run lint
  ```

- **Format** (Prettier):

  ```bash
  cd frontend
  npm run format
  ```

### Environment & config notes

- Vite uses `import.meta.env` for environment variables (see `vite.config.ts` and `tsconfig` files).
- For Docker-based deployments, the frontend image is built with arguments from `docker/services/frontend/Dockerfile` and shares its build output with Nginx.
- When integrating with the backend, keep the API base URL and auth flows in sync with the Go API routes and config.
