CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE administrators (
  id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),

  username        VARCHAR(100) NOT NULL UNIQUE,
  password_hash   VARCHAR(255) NOT NULL,

  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_administrators_username ON administrators (username);
