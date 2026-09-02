CREATE TABLE loans (
  id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),

  student_id      UUID        NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
  book_id         UUID        NOT NULL REFERENCES books(id)    ON DELETE RESTRICT,

  loan_date       TIMESTAMPTZ NOT NULL,
  due_date        TIMESTAMPTZ NOT NULL,
  return_date     TIMESTAMPTZ,
  status          VARCHAR(20) NOT NULL DEFAULT 'ACTIVE'
                  CHECK (status IN ('ACTIVE', 'RETURNED')),
  was_late        BOOLEAN,

  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_loans_student_id ON loans (student_id);
CREATE INDEX idx_loans_book_id ON loans (book_id);
CREATE INDEX idx_loans_status_due_date ON loans (status, due_date);
