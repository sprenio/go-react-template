---
trigger: model_decision
description: Use this rule when you need to verify sent emails, retrieve confirmation links, check password reset tokens etc. Activate this whenever a task involves reading an email sent by application.
---

# MailHog Email Access Protocol

## 1. Logic for MailHog URL
You are not allowed to guess the MailHog URL. You must derive it dynamically based on the current Application URL being tested using the following transformation rule:
- **Rule:** Replace the application subdomain (e.g., `app`, `project`, `api`) with `mailhog`, preserving the `www` prefix if present.
- **Logic:** 
  - `[app-subdomain].[domain].[tld]` → `mailhog.[domain].[tld]`
  - `www.[app-subdomain].[domain].[tld]` → `www.mailhog.[domain].[tld]`

**Examples:**
- `www.project.localhost` → `www.mailhog.localhost`
- `app.mydomain.com` → `mailhog.mydomain.com`
- `api.project.dev` → `mailhog.project.dev`

## 2. When to Access MailHog
Trigger this protocol if:
- You just performed an action that triggers an email (registration, password reset, alerts).
- You need a **Verification Code** or **Confirmation Link** to proceed with a test.
- You need to verify if the email body/formatting is correct.

## 3. Workflow for Agent
1. **Identify the App URL:** Determine the base URL you are currently using for the application.
2. **Calculate MailHog URL:** Apply the transformation rule above.
3. **Open Browser:** Use your browser tool to navigate to the calculated MailHog URL.
4. **Retrieve Data:** Find the latest email, extract the necessary information (e.g., click a link or copy a token), and return to the application flow.

## 4. Interaction Guidelines
- If the MailHog UI is empty, wait 5 seconds and refresh the browser before reporting a failure.
- Always report the specific link or code you found in the email to the user for transparency.
