---
trigger: model_decision
description: Use this rule when adding, modifying, or deleting translations in the frontend. It defines the folder structure for language files (pl, en, de, uk), JSON format requirements, and the logic for grouping keys into modules.
---

# Frontend Internationalization (i18n) Protocol

## 1. Directory Structure
All translation files are located in: `frontend/src/locales/`
Languages are organized into subdirectories (ISO codes):
- `/pl/`, `/en/`, `/de/`, `/uk/`

## 2. JSON File Organization
The filename (e.g., `login.json`) determines the **Module Group** in the application.
- **`base.json`**: Shared keys, `apiCodes`, `loader`, and generic UI (e.g., "Cancel").
- **Module files**: Specific to a feature (e.g., `login.json`, `register.json`, `dashboard.json`, `settings.json`).

## 3. Merging Logic
The loader script `frontend/src/i18n/index.ts` merges all JSONs.
- Access global keys: `t('apiCodes.1101')` or `t('generic.cancel')`
- Access module keys: `t('login.submit')` (prefixed by filename)

## 4. Mandatory Coding Standards
- **ZERO Hardcoded Strings:** You are STRICTLY FORBIDDEN from putting raw text in React components. Always use the `t()` function from your translation hook.
- **Full Synchronization:** When adding a new key, you MUST add it to the corresponding JSON file in **EVERY** language directory (`pl`, `en`, `de`, `uk`).
- **Naming:** Use `snake_case`. Use prefixes: `lbl_` (label), `ph_` (placeholder), `btn_` (button), `err_` (error).

## 5. Workflow for New Translations
1. Identify if it's a `base` or `module` string.
2. Locate the file in `frontend/src/locales/pl/` (or `en/`).
3. Add the key and value.
4. **Immediately** replicate that key in all other language folders (`pl/`, `en/`, `de/`, `uk/`) with translated values.

## 6. Adding a New Language
If a new language (e.g., `fr`) is requested:
1. **Create Directory:** Create `frontend/src/locales/fr/`.
2. **Copy All Files:** Every single JSON file from `en/` (or `pl/`) must be copied to `fr/` and translated. No file can be missing.
3. **Update Loader:** You MUST modify `frontend/src/i18n/index.ts` to import the new language and add it to the i18next configuration.
