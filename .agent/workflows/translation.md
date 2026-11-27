---
description: Translation Synchronization Workflow (Frontend + Backend)
---

# Translation Synchronization Workflow (Frontend + Backend)

## Goal

Using `pl` (Polish) as the source of truth:

* Synchronize missing / outdated translations in all other languages (`en`, `de`, `uk`, etc.)
* Ensure identical file structure and translation keys
* Generate translations in proper product context (not literal word‑by‑word)

---

# 1. FRONTEND

## Source Structure (Reference)

```
./frontend/src/locales/pl/{module}.json
```

## Target Structure

```
./frontend/src/locales/{lang_code}/{module}.json
```

---

## Rules

1. `pl` is the single source of truth.
2. Every language must:

   * Have exactly the same `{module}.json` files as `pl`
   * Contain exactly the same translation keys inside each file
3. If a module file does not exist → create it.
4. If a key exists in `pl` but not in target language → add it.
5. If a key was removed from `pl` → remove it from other languages.
6. Do NOT modify existing translations unless the Polish source meaning changed.
7. Preserve:

   * JSON structure
   * Key names
   * Placeholders (`{{variable}}`, `{count}`, etc.)
   * Formatting

---

## Frontend Algorithm

For each language ≠ `pl`:

1. Iterate over all `{module}.json` files in `pl`
2. If module file does not exist → create and translate all keys
3. If module file exists:

   * Compare keys with `pl`
   * Add missing keys
   * Remove extra keys
4. Keep structure identical to `pl`

Translations must reflect:

* UI context (e.g., `login.json` → authentication screen)
* Product tone (consistent with Polish)
* Natural UX wording

---

# 2. BACKEND

## Source Structure (Reference)

```
./backend/locale/translates/pl.json
```

Each entry structure:

```json
{
  "id": "welcome.confirm_email",
  "translation": "Potwierdź adres e-mail"
}
```

## Target Structure

```
./backend/locale/translates/{lang_code}.json
```

---

## Rules

1. `pl.json` is the single source of truth.
2. Every language must contain exactly the same `id` values.
3. Maintain the same order as in `pl.json`.
4. If an `id`:

   * Exists in `pl` but not in target → add it
   * Was removed from `pl` → remove it from target
5. Preserve:

   * HTML (e.g., `<br/>`)
   * Backend interpolation syntax (`{{ .AppName }}`)
   * Variables inside placeholders
6. Never translate placeholder variable names.

---

## Backend Algorithm

For each language ≠ `pl`:

1. Load `pl.json`
2. Load `{lang}.json`
3. For each `id` in `pl.json`:

   * If missing → add translated entry
4. Remove entries not present in `pl.json`
5. Ensure identical ordering

---

# 3. Translation Guidelines

* Do NOT translate:

  * Keys
  * Placeholder variable names
  * HTML tags
* Keep translations concise and product‑appropriate
* Preserve tone consistency
* Avoid unnecessary expansion of text length
* Do not add comments to JSON files

---

# 4. Execution Modes

## Full Synchronization

Trigger phrase:

> sync translations

Actions:

* All frontend modules
* Backend
* All languages except `pl`

---

## Frontend Module Only

Trigger phrase:

> sync translations frontend {module}

Actions:

* Only specified module
* All languages except `pl`

---

## Backend Only

Trigger phrase:

> sync translations backend

---

# 5. Completion Criteria

Process is complete when:

* Every language mirrors `pl` structure exactly
* No missing keys
* No extra keys
* Valid JSON syntax in all files
