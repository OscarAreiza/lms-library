CREATE TABLE students (
  id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),

  full_name         VARCHAR(200) NOT NULL,
  document_id       VARCHAR(50)  NOT NULL UNIQUE,
  email             VARCHAR(255) NOT NULL,
  phone             VARCHAR(30),

  suspended_until   TIMESTAMPTZ,
  deactivated_at    TIMESTAMPTZ,

  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_students_document_id ON students (document_id);
CREATE INDEX idx_students_active ON students (deactivated_at) WHERE deactivated_at IS NULL;
CREATE INDEX idx_students_suspended_until ON students (suspended_until) WHERE suspended_until IS NOT NULL;
