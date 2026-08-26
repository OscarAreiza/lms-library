package postgres

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OscarAreiza/lms-library/backend/internal/domain/circulation"
)

// LoanRepository implements circulation.LoanRepository against PostgreSQL —
// the only code allowed to touch the `loans` table
// (library-docs/09-microservices/service-boundary-rules.md).
type LoanRepository struct {
	db *pgxpool.Pool
}

func NewLoanRepository(db *pgxpool.Pool) *LoanRepository {
	return &LoanRepository{db: db}
}

const loanColumns = `id, student_id, book_id, loan_date, due_date, return_date, status, was_late, created_at, updated_at`

func scanLoan(row pgx.Row) (*circulation.Loan, error) {
	var l circulation.Loan
	if err := row.Scan(&l.ID, &l.StudentID, &l.BookID, &l.LoanDate, &l.DueDate, &l.ReturnDate, &l.Status, &l.WasLate, &l.CreatedAt, &l.UpdatedAt); err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LoanRepository) FindByID(ctx context.Context, id string) (*circulation.Loan, error) {
	query := `SELECT ` + loanColumns + ` FROM loans WHERE id = $1`
	l, err := scanLoan(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, circulation.ErrLoanNotFound
	}
	return l, err
}

func (r *LoanRepository) CountActiveByStudent(ctx context.Context, studentID string) (int, error) {
	const query = `SELECT count(*) FROM loans WHERE student_id = $1 AND status = 'ACTIVE'`
	var count int
	err := r.db.QueryRow(ctx, query, studentID).Scan(&count)
	return count, err
}

// Search builds a filtered query for HU-07/HU-08 (loan history, overdue loans).
// Filters are applied incrementally since most callers only need a subset.
func (r *LoanRepository) Search(ctx context.Context, status string, overdueOnly bool, studentID, bookID string, page, limit int) ([]*circulation.Loan, int, error) {
	offset := (page - 1) * limit
	query := `SELECT ` + loanColumns + ` FROM loans WHERE 1=1`
	countQuery := `SELECT count(*) FROM loans WHERE 1=1`
	var args []any

	addFilter := func(clause string, value any) {
		args = append(args, value)
		placeholder := " AND " + clause + " = $" + strconv.Itoa(len(args))
		query += placeholder
		countQuery += placeholder
	}

	if status != "" {
		addFilter("status", status)
	}
	if overdueOnly {
		query += " AND status = 'ACTIVE' AND due_date < now()"
		countQuery += " AND status = 'ACTIVE' AND due_date < now()"
	}
	if studentID != "" {
		addFilter("student_id", studentID)
	}
	if bookID != "" {
		addFilter("book_id", bookID)
	}

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitArgs := append(append([]any{}, args...), limit, offset)
	query += " ORDER BY loan_date DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)

	rows, err := r.db.Query(ctx, query, limitArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var loans []*circulation.Loan
	for rows.Next() {
		l, err := scanLoan(rows)
		if err != nil {
			return nil, 0, err
		}
		loans = append(loans, l)
	}

	return loans, total, nil
}

func (r *LoanRepository) Save(ctx context.Context, l *circulation.Loan) error {
	const query = `
		INSERT INTO loans (id, student_id, book_id, loan_date, due_date, return_date, status, was_late, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			return_date = EXCLUDED.return_date,
			status = EXCLUDED.status,
			was_late = EXCLUDED.was_late,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.Exec(ctx, query, l.ID, l.StudentID, l.BookID, l.LoanDate, l.DueDate, l.ReturnDate, l.Status, l.WasLate, l.CreatedAt, l.UpdatedAt)
	return err
}
