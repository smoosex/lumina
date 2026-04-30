# Lumina

[中文](./README.zh-CN.md) | English

Lumina is my personal digital garden: a portfolio homepage plus a reading-notes section for publishing book reflections written locally in Markdown.

<p align="center">
  <img src="frontend/previews/preview1.png" width="45%" />
  <img src="frontend/previews/preview2.png" width="45%" />
  <img src="frontend/previews/preview3.png" width="45%" />
  <img src="frontend/previews/preview4.png" width="45%" />
  <img src="frontend/previews/preview5.png" width="45%" />
  <img src="frontend/previews/preview6.png" width="45%" />
  <img src="frontend/previews/preview7.png" width="45%" />
  <img src="frontend/previews/preview8.png" width="45%" />
</p>

The project is now a monorepo:

```text
lumina/
├── frontend/   # Vue 3 portfolio and reading notes UI
└── backend/    # Go API for reading notes
```

## Features

- Portfolio homepage with animated hero, skills, selected works, and contact links.
- Reading notes list and article detail pages.
- Markdown upload workflow protected by signed requests instead of a login system.
- Book metadata fetched by ISBN during upload.
- Local book covers saved and served by the backend.

## Tech Stack

- Frontend: Vue 3, Vite, TypeScript, Tailwind CSS v4, Pinia, vue-i18n, GSAP, Three.js.
- Backend: Go, Fx, Gin, Gorm.
- Database: Supabase PostgreSQL.
- Runtime/development: Bun for frontend, Go toolchain for backend.

## Quick Start

Install frontend dependencies:

```bash
cd frontend
bun install
bun dev
```

Start the backend:

```bash
cd backend
cp .env.example .env
go run ./cmd/api
```

Run backend tests:

```bash
cd backend
go test ./...
```

Build the frontend:

```bash
cd frontend
bun run build
```

## Configuration

Backend configuration lives in `backend/.env`:

```text
SERVER_ADDR=127.0.0.1:8080
DATABASE_DSN=postgres://USER:PASSWORD@HOST:6543/postgres?sslmode=require
UPLOAD_SIGNATURE_SECRET=change-this-long-random-secret
UPLOAD_API_BASE_URL=http://127.0.0.1:8080
```

## Reading Notes

A reading note is a local Markdown file with frontmatter:

```markdown
---
slug: psychology-of-money
title: "读完《金钱心理学》，我开始少一点轻易评判别人"
isbn: "9787513941242"
status: published
tags:
  - 金钱
  - 心理学
---

# 读完《金钱心理学》，我开始少一点轻易评判别人

正文内容。
```

Upload a note:

```bash
cd backend
go run ./cmd/upload-note -file ../notes/psychology-of-money.md
```

The backend fetches book metadata by ISBN during upload, renders the Markdown, stores it in PostgreSQL, and downloads the cover once.

## Documentation

- [Frontend README](frontend/README.md)
- [Frontend README 中文版](frontend/README.zh-CN.md)
