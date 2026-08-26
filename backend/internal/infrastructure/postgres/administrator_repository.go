package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OscarAreiza/lms-library/backend/internal/domain/access"
)

// AdministratorRepository implements access.AdministratorRepository against
// PostgreSQL — this is the only code allowed to touch the `administrators`
// table (library-docs/09-microservices/service-boundary-rules.md).
type AdministratorRepository struct {
	db *pgxpool.Pool
}

func NewAdministratorRepository(db *pgxpool.Pool) *AdministratorRepository {
	return &AdministratorRepository{db: db}
}

func (r *AdministratorRepository) FindByUsername(ctx context.Context, username string) (*access.Administrator, error) {
	const query = `
		SELECT id, username, password_hash, created_at, updated_at
		FROM administrators
		WHERE username = $1`

	var a access.Administrator
	err := r.db.QueryRow(ctx, query, username).Scan(&a.ID, &a.Username, &a.PasswordHash, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, access.ErrAdministratorNotFound
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AdministratorRepository) Save(ctx context.Context, a *access.Administrator) error {
	const query = `
		INSERT INTO administrators (id, username, password_hash, created_at, updated_at)
		VALUES (COALESCE(NULLIF($1, ''), gen_random_uuid()::text)::uuid, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			password_hash = EXCLUDED.password_hash,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.Exec(ctx, query, a.ID, a.Username, a.PasswordHash, a.CreatedAt, a.UpdatedAt)
	return err
}
