---
trigger: always_on
---

# Overview: Application Template (go-react-template)
Core Concept: Full-stack Go + React + TypeScript boilerplate for web applications with authentication, user management, asynchronous processing, and dockerized local development.

# System Architecture

## Tech Stack

### Backend
* **golang 1.26** - REST API with JWT-based authentication, queue consumers, cron worker

### Frontend
* **react 19.2**
* **typescript 5.9.3**
* **tailwindcss 4.1**
* **vite 7.2.4**
* **vitest 4.0.16**

### Datastore & Caching
* **MariaDB 13.0** - primary transactional DB.
* **Redis 8.8** - caching, distributed locks, rate limiting.
* **PHP 8.3** - script for data migration.

### Messaging, Async
* **RabbitMQ 4.1** - job for consumers.

### Hosting and Network
* **Nginx 1.29** - web server.

## Project repository
Key repository directories and their roles:
* **.agent**: Agent configuration, rules (`.agent/rules`), workflows (`.agent/workflows`), and helper scripts (`.agent/scripts`).
* **backend**: Go REST API, config, services, middleware, queue consumers, cron worker.
* **db_migrations**: SQL migrations for MariaDB schema.
* **docker**: Dockerfiles, compose files, Nginx config, helper script (`docker-control.sh`).
* **frontend**: React + TS SPA with routes, hooks, providers, and Tailwind styling.
* **storage**: Docker volumes for database, RabbitMQ, Redis, and Nginx logs.

Note on helper scripts:
Agent and developer helper scripts are located in `.agent/scripts`. If you create a new helper script, place it in `.agent/scripts`.


------------------

# Chat rules.
* **ALWAYS** reply in Polish in chat, unless the direct instructions state otherwise.
* If you create a plan or task, **ALWAYS** use Polish there.


------------------

# Reguły projektu i niespójności (Project Rules & Inconsistencies)
* Projekt jest w fazie dynamicznego rozwoju i koncepcje mogą ulegać gwałtownym zmianom bez natychmiastowej aktualizacji reguł. W związku z tym niektóre wytyczne mogą być nieaktualne, mylące lub niespójne ze stanem faktycznym.
* **ZAKAZ DOMYŚLANIA SIĘ I SAMOWOLNEGO ROZWIĄZYWANIA WĄTPLIWOŚCI:** Jeśli cokolwiek w regułach wydaje się nieaktualne, problematyczne, sprzeczne z kodem lub wewnętrznie niespójne – **pod żadnym pozorem nie próbuj samodzielnie zgadywać intencji ani "domyślać się" rozwiązania na własną rękę**.
* **OBOWIĄZEK NATYCHMIASTOWEGO ZGŁOSZENIA:** Zawsze od razu poinformuj użytkownika o zauważonej niespójności, wątpliwości lub nieaktualnej regule i poproś o doprecyzowanie, zanim podejmiesz jakiekolwiek działania.


------------------

# **Key documentation** - Documentation of the most important mechanics and classes of the system

Please read carefully when reviewing the documentation provided in the specified MD file. All the most important system mechanics that may be useful to you when analysing or creating code are documented.

* `./development.md` - Rules of Development and Code Writing and Code Review Best Practices.
* `./tools.md` - The list of tools you can use.
* `./docker.md` - Docker Infrastructure & Management Rules.
* `./frontend.md` - Frontend Development Standards
* `./backend.md` - Backend Development Standards
* `./ratelimit.md` - Rate Limiting in Backend API
