---
trigger: glob
description: General information about frontend. Read it ALWAYS AND ONLY if the question or task is related to the frontend
globs: frontend/**
---

Use this rule when working on the React frontend. It defines the project structure (located in frontend/), coding standards, Tailwind CSS patterns with theme variables, and component organization.


# Frontend Development Standards (React + TypeScript)

## 1. Project Location & Context
- **Root Directory:** `frontend`
- **Source Code:** All development happens inside `frontend/src`.
- **Static Assets:** Located in `./public` (logos, favicons).

## 2. Core Architecture (`frontend/src/`)
- **main.tsx:** Entry point. Wraps the app in `StrictMode`, `AppProviders` (Auth, Lang, Config), and `CustomToaster`.
- **App.tsx:** Main router. Controls flow between Protected and Public routes based on Auth status.

### Directory Structure:
- `api/`: API communication logic (e.g., `apiClient`, `apiCodes`).
- `assets/`: Component-specific assets (logos, icons).
- `components/`: Reusable UI elements (Button, Card, Form, Modal, etc.).
- `domain/`: Global TypeScript types and interfaces (e.g., `config.ts`, `user.ts`).
- `helpers/`: Generic utility functions.
- `hooks/`: Custom React hooks and mutations (e.g., `useLogin`, `useRegister`).
- `layouts/`: Base layouts (`AuthenticatedLayout`, `PublicLayout`) and common UI (Header, Footer).
- `locales/`: i18n translations. Refer to `frontend-internationalization.md` for details.
- `pages/`: Page-level components organized by route.
- `providers/`: Context providers (Auth, Theme, Config).
- `routes/`: Logic for `PublicRoute` and `ProtectedRoute`.
- `styles/`: Global CSS and Theme definitions (`theme.base.css`, `theme.light.css`, `theme.dark.css`, `globals.css`, `index.css`).
- `test/`: Test suites.

## 3. Coding Standards & Rules

### Component Exports & Imports
- **Barrel Exports:** Every component directory MUST contain an `index.ts` file that handles the export.
  - *Example (`index.ts`):* `export { default as LanguageSelect } from './LanguageSelect';`
- **Clean Imports:** NEVER import directly from the component file. ALWAYS import from the directory level.
  - ❌ **BAD:** `import { LanguageSelect } from '@/components/LanguageSelect/LanguageSelect.tsx';`
  - ✅ **GOOD:** `import { LanguageSelect } from '@/components/LanguageSelect';`
- **Exception:** You may use relative direct imports (e.g., `./SubComponent`) ONLY when importing something from the same directory (the same component space).

### Component Placement
- **Global Components:** Small, reusable wrappers go to `src/components`.
- **Feature-Specific Components:** If a component (e.g., `LoginForm`) is used only on one page, place it in that page's directory (e.g., `src/pages/Login/LoginForm.tsx`) and provide an `index.ts` for it.

### Styling (Tailwind CSS & Themes)
- **NO Inline Styles:** Use Tailwind classes only.
- **NO Hardcoded Values:** Never hardcode colors, spacing, or font sizes.
- **Theme Variables:** Always use CSS variables from the theme for padding, margins, colors, and text sizes.
  - ❌ **BAD:** `className="w-full text-lg px-6 py-3 text-white"`
  - ✅ **GOOD:** `className="w-full px-[var(--space-lg)] py-[var(--space-md)] text-[length:var(--text-lg)]"`
- **Dark/Light Mode:** When styling, ensure the element looks correct in both themes. If a new color is needed, add a variable to `theme.light.css` and `theme.dark.css`.

### Clean Code & Logs
- **Comments:** Use ONLY when absolutely necessary and ONLY in **English**.
- **Logging:** Use `console.warn` or `console.error`. **NO `console.log` or debug logs** in production code.
- **Linting:** Code must be compliant with ESLint rules (no warnings/errors allowed).
