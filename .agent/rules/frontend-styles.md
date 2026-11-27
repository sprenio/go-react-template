---
trigger: model_decision
description: This rule must be applied whenever the agent is creating, modifying, or reviewing any UI layout, styling, or visual aspects of the frontend.
---

### General Principles
- All newly created or modified UI elements must remain visually consistent with the rest of the application.
- Reuse existing patterns, components, and styling conventions whenever possible.
- Avoid introducing new visual styles if an equivalent already exists in the project.

### Styling System
- The project uses **TailwindCSS (version: 4.1.18)** as the primary styling framework.
- Prefer Tailwind utility classes over custom CSS whenever possible.
- Use existing helper classes defined in the project before introducing new ones.

### Theme System
- The application uses a theme-based design system.
- Theme variables are defined in:
  - `frontend/src/styles/theme.base.css` – shared base variables
  - `frontend/src/styles/theme.dark.css` – dark theme variables
  - `frontend/src/styles/theme.light.css` – light theme variables

### CSS Structure
- Additional styles are organized as follows:
  - `frontend/src/styles/globals.css` – helper classes (often grouping Tailwind utilities)
  - `frontend/src/styles/index.css` – main entry point importing all styles

### Colors and Design Tokens
- **Never use hardcoded colors.**
- Always use theme variables defined in the theme files.
- Ensure all colors adapt correctly to both light and dark themes.

### Spacing and Layout
- Prefer using predefined spacing values (e.g. padding, margin) if available.
- Reuse existing spacing conventions to maintain visual rhythm and consistency.
- Avoid arbitrary values unless absolutely necessary.

### Maintainability
- Do not duplicate styles unnecessarily.
- If a pattern repeats, consider extracting it into a reusable helper or component.
- Keep styling predictable, scalable, and aligned with the existing system.

### When in Doubt
- Match the closest existing UI pattern in the codebase.
- Prioritize consistency over creativity.
