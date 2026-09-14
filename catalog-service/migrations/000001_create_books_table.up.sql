CREATE TABLE books (
  id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),

  title              VARCHAR(300) NOT NULL,
  author             VARCHAR(200) NOT NULL,
  isbn               VARCHAR(20)  NOT NULL UNIQUE,
  category           VARCHAR(100) NOT NULL,
  year               INTEGER      NOT NULL,

  total_copies       INTEGER      NOT NULL CHECK (total_copies >= 1),
  available_copies   INTEGER      NOT NULL CHECK (available_copies >= 0 AND available_copies <= total_copies),

  created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_books_isbn ON books (isbn);
CREATE INDEX idx_books_title ON books (title);
CREATE INDEX idx_books_author ON books (author);
CREATE INDEX idx_books_category ON books (category);
