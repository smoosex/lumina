create extension if not exists pgcrypto;

create table if not exists reading_notes (
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

create table if not exists reading_tags (
  id uuid primary key default gen_random_uuid(),
  name text not null unique
);

create table if not exists reading_note_tags (
  note_id uuid not null references reading_notes(id) on delete cascade,
  tag_id uuid not null references reading_tags(id) on delete cascade,
  primary key (note_id, tag_id)
);

create index if not exists idx_reading_notes_status_published_at
  on reading_notes (status, published_at desc);

create index if not exists idx_reading_notes_slug
  on reading_notes (slug);

alter table reading_notes enable row level security;
alter table reading_tags enable row level security;
alter table reading_note_tags enable row level security;
