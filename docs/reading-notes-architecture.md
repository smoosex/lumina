# Reading Notes Architecture

## Goal

Add a reading notes column to Lumina for publishing book reflections written locally as Markdown.

The feature is not a CMS and does not need a login system. The intended workflow is:

1. Write a Markdown file locally.
2. Upload the file with a local script or CLI.
3. The Go backend parses, validates, renders, and stores the note.
4. The frontend only consumes public read APIs.

## Repository Layout

Move the current Vue project into `frontend/` and add a sibling Go backend.

```text
lumina/
├── docs/
│   └── reading-notes-architecture.md
├── frontend/                # current Vue project moves here
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── vite.config.ts
├── backend/
│   ├── cmd/api/main.go
│   ├── .env.example         # committed backend example
│   ├── .env                 # server/local only, not committed
│   ├── internal/
│   │   ├── app/              # fx module wiring
│   │   ├── config/           # .env loader and typed config
│   │   ├── db/               # gorm + postgres
│   │   ├── http/             # gin router and middleware
│   │   ├── notes/            # domain, service, repository
│   │   ├── markdown/         # frontmatter, markdown render, sanitize
│   │   └── uploadauth/       # signed upload middleware
│   ├── migrations/
│   └── go.mod
```

## Tech Choices

- Frontend: current Vue 3 + Vite app.
- Backend: Go + Fx + Gin + Gorm.
- Database: Supabase PostgreSQL.
- Backend configuration: `.env` file loaded by the Go process.
- Frontend configuration: no environment variables required for same-origin Nginx deployment.
- Content source: local Markdown files with frontmatter.

## Backend Configuration

Use `backend/.env` for backend configuration. Commit only `backend/.env.example`; never commit real `.env` files.

```text
SERVER_ADDR=127.0.0.1:8080
DATABASE_DSN=postgres://USER:PASSWORD@HOST:6543/postgres?sslmode=require
UPLOAD_SIGNATURE_SECRET=change-this-long-random-secret
UPLOAD_API_BASE_URL=https://www.smoose.cn/lumina
```

The repository should ignore real env files:

```text
.env
.env.*
!.env.example
!**/.env.example
```

Copy `backend/.env.example` to `backend/.env`, then replace placeholder values. The frontend calls same-origin `/lumina/api/*` and `/lumina/covers/*`, so it does not need its own `.env` for the Nginx deployment.

## Secure Upload Without Login

Use signed upload requests instead of a login system.

The local uploader and backend share one long random secret from `.env`. Each upload request includes:

```text
X-Lumina-Timestamp: unix seconds
X-Lumina-Body-SHA256: sha256 hex of request body
X-Lumina-Signature: hmac-sha256(secret, method + "\n" + path + "\n" + timestamp + "\n" + body_sha256)
```

Backend validation:

- Reject missing signature headers.
- Reject timestamps outside the configured skew window.
- Recalculate the body hash and HMAC signature.
- Use constant-time comparison.
- Enforce Markdown file size limit.

This keeps the public API readable by everyone while making write access possible only from your local uploader. A plain Bearer token would also work, but signed requests are better because replay is limited and accidental leakage is less catastrophic. Static bearer token is the lazy cousin who still somehow gets invited to production.

## Markdown Format

Each note is one Markdown file with frontmatter.

```markdown
---
slug: thinking-in-systems
title: "读《Thinking in Systems》之后"
isbn: "9781603580557"
status: published
rating: 5
read_at: 2026-04-20
published_at: 2026-04-29T10:00:00+08:00
tags:
  - Thinking
  - Systems
excerpt: "这本书让我重新理解了反馈、延迟和结构对结果的影响。"
---

# 读《Thinking in Systems》之后

这里写完整读后感正文。
```

Required fields:

- `slug`
- `title`
- `isbn`
- `status`
- Markdown body

Recommended fields:

- `tags`
- `excerpt`
- `read_at`
- `published_at`
- `rating`

## Database Schema

```sql
create table reading_notes (
  id uuid primary key default gen_random_uuid(),
  slug text not null unique,
  title text not null,
  isbn text,
  douban_id text,
  book_title text not null,
  author text,
  translator text,
  publisher text,
  producer text,
  book_pub_date text,
  book_pages int,
  cover_url text,
  douban_url text,
  excerpt text,
  content_md text not null,
  content_html text,
  status text not null default 'draft',
  rating int,
  read_at date,
  published_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table reading_tags (
  id uuid primary key default gen_random_uuid(),
  name text not null unique
);

create table reading_note_tags (
  note_id uuid not null references reading_notes(id) on delete cascade,
  tag_id uuid not null references reading_tags(id) on delete cascade,
  primary key (note_id, tag_id)
);

create index idx_reading_notes_status_published_at
  on reading_notes (status, published_at desc);

create index idx_reading_notes_slug
  on reading_notes (slug);

create index idx_reading_notes_isbn
  on reading_notes (isbn);
```

Gorm models should map to these tables directly. Keep migrations as SQL files instead of relying only on `AutoMigrate`, because Supabase is real infrastructure and deserves fewer surprises.

## API Design

Public display APIs:

```text
GET /api/reading-notes
GET /api/reading-notes/:slug
GET /api/reading-tags
```

Private upload API:

```text
PUT /api/admin/reading-notes/:slug
```

The upload endpoint accepts either:

```text
Content-Type: text/markdown
```

or:

```text
Content-Type: multipart/form-data
file=<markdown file>
```

The server parses the frontmatter, checks that path `:slug` matches frontmatter `slug`, fetches Douban book metadata by `isbn`, renders Markdown, sanitizes HTML, and upserts the database row. Douban is only requested during upload; the frontend reads persisted metadata from PostgreSQL.

## Response Shapes

List:

```json
{
  "items": [
    {
      "slug": "thinking-in-systems",
      "title": "读《Thinking in Systems》之后",
      "isbn": "9781603580557",
      "doubanId": "3793669",
      "bookTitle": "Thinking in Systems",
      "author": "Donella H. Meadows",
      "publisher": "Chelsea Green Publishing",
      "bookPubDate": "2008-12-3",
      "coverUrl": "https://img9.doubanio.com/view/subject/l/public/...",
      "doubanUrl": "https://book.douban.com/subject/3793669/",
      "excerpt": "这本书让我重新理解了反馈、延迟和结构对结果的影响。",
      "tags": ["Thinking", "Systems"],
      "rating": 5,
      "readAt": "2026-04-20",
      "publishedAt": "2026-04-29T10:00:00+08:00"
    }
  ]
}
```

Detail:

```json
{
  "slug": "thinking-in-systems",
  "title": "读《Thinking in Systems》之后",
  "isbn": "9781603580557",
  "doubanId": "3793669",
  "bookTitle": "Thinking in Systems",
  "author": "Donella H. Meadows",
  "publisher": "Chelsea Green Publishing",
  "bookPubDate": "2008-12-3",
  "coverUrl": "https://img9.doubanio.com/view/subject/l/public/...",
  "doubanUrl": "https://book.douban.com/subject/3793669/",
  "excerpt": "这本书让我重新理解了反馈、延迟和结构对结果的影响。",
  "contentHtml": "<article>...</article>",
  "tags": ["Thinking", "Systems"],
  "rating": 5,
  "readAt": "2026-04-20",
  "publishedAt": "2026-04-29T10:00:00+08:00"
}
```

## Backend Module Design

Fx modules:

```text
config.Module
db.Module
bookmeta.Module
markdown.Module
uploadauth.Module
notes.Module
http.Module
```

Responsibilities:

- `config`: load `.env` once, validate required keys, and provide typed config.
- `db`: open Supabase PostgreSQL connection through Gorm.
- `bookmeta`: fetch and parse Douban book metadata by ISBN during upload.
- `markdown`: parse frontmatter, render Markdown with `goldmark`, sanitize with `bluemonday`.
- `uploadauth`: Gin middleware for signed upload routes.
- `notes`: repository and service for list/detail/upsert.
- `http`: Gin router, health check, API routes.

Suggested libraries:

```text
go.uber.org/fx
github.com/gin-gonic/gin
gorm.io/gorm
gorm.io/driver/postgres
github.com/yuin/goldmark
github.com/microcosm-cc/bluemonday
github.com/joho/godotenv
gopkg.in/yaml.v3
```

## Frontend Integration

The frontend is served under `/lumina/`. It requests `/lumina/api/*` and `/lumina/covers/*`; in production, Nginx rewrites those paths to the Go backend. In local development, Vite proxies them to `http://127.0.0.1:8080`.

Frontend modules:

```text
frontend/src/features/reading-notes/
├── api.ts
└── types.ts

frontend/src/views/ReadingNotesView.vue
frontend/src/views/ReadingNoteView.vue
```

Routes are handled by a lightweight path switch in `App.vue` for now, without adding `vue-router`:

```text
/notes
/notes/:slug
```

The home page only links to reading notes. The reading notes list lives on `/lumina/notes`, and detail links use `/lumina/notes/:slug`, so the self-hosted frontend server needs history fallback to `frontend/dist/index.html`.

## Upload CLI Flow

The local upload command should:

1. Read Markdown file.
2. Parse frontmatter locally for quick validation.
3. Compute body SHA256.
4. Sign request with the local secret.
5. Send `PUT /lumina/api/admin/reading-notes/:slug`.
6. Print the published URL.

Example:

```bash
cd backend
go run ./cmd/upload-note -file ../notes/thinking-in-systems.md
```

The CLI reads `UPLOAD_SIGNATURE_SECRET` and `UPLOAD_API_BASE_URL` from `backend/.env`. `-base-url` can override the target API when testing locally.

## Implementation Order

1. Move the current Vue app into `frontend/` and adjust scripts/paths.
2. Create `backend/` skeleton with Fx, Gin, `.env` config loader, and health check.
3. Add SQL migrations and Gorm models.
4. Implement Markdown parser/render/sanitizer.
5. Implement signed upload middleware.
6. Implement note upsert endpoint.
7. Implement public list/detail APIs.
8. Add frontend API client that requests same-origin `/lumina/api/*`.
9. Build standalone Reading Notes list page.
10. Add article-style reading note page.
11. Add local upload CLI.

## Non-goals

- No user login.
- No web admin panel for now.
- No browser-based Markdown editor.
- No direct Supabase access from the frontend.
- No public write APIs.
