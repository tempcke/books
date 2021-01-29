CREATE TABLE IF NOT EXISTS books (
  id         VARCHAR(36)  PRIMARY KEY,
  title      VARCHAR(128),
  author     VARCHAR(128),
  pubdate    date,
  rating     INT,
  status     VARCHAR(16),
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);
