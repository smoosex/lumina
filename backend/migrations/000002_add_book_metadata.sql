alter table reading_notes
  add column if not exists isbn text,
  add column if not exists douban_id text,
  add column if not exists translator text,
  add column if not exists publisher text,
  add column if not exists producer text,
  add column if not exists book_pub_date text,
  add column if not exists book_pages int,
  add column if not exists cover_url text,
  add column if not exists douban_url text;

create index if not exists idx_reading_notes_isbn
  on reading_notes (isbn);
