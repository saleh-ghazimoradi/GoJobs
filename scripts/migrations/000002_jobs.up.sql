CREATE TABLE IF NOT EXISTS jobs (
  id bigserial PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  location TEXT NOT NULL,
  company TEXT NOT NULL,
  salary TEXT NOT NULL,
  user_id bigint NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES users(id)
);