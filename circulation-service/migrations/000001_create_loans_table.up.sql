-- student_id/book_id are plain UUID columns, not foreign keys — Student and
-- Book now live in membership-service's and catalog-service's own databases
-- (library-docs/09-microservices/service-boundary-rules.md). Referential
-- integrity across services is enforced at the application layer (HTTP
-- calls in LoanRegistrationService), not by the database.
CREATE TABLE loans (
  id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),

  student_id      UUID        NOT NULL,
  book_id         UUID        NOT NULL,

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
