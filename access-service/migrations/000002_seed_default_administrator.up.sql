-- Local/dev seed only — HU-01. Do not run this against a real production
-- database with this fixed password; provision production admins out-of-band
-- (library-docs/02-domain/entities-and-rules.md: Administrator is never
-- self-registered).
INSERT INTO administrators (username, password_hash)
VALUES ('admin@lms.com', '$2a$12$GUrCkw4mvNJUmdN6pPb/5.7DKtye8C0vTpcO1d43UvixvZgWpBLKq')
ON CONFLICT (username) DO NOTHING;
