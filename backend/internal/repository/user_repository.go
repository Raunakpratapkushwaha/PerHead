package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	UpdateQRCode(ctx context.Context, userID uint64, qrURL string) error
}

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) CreateUser(ctx context.Context, u *model.User) error {
	// 1. Changed ? to $1, $2, etc.
	// 2. Added RETURNING id for PostgreSQL
	query := `INSERT INTO users (name, email, phone, password_hash) VALUES ($1, $2, $3, $4) RETURNING id`

	// PostgreSQL doesn't support LastInsertId() on Exec, so we QueryRow and Scan the returned ID
	err := r.db.QueryRowContext(ctx, query, u.Name, u.Email, u.Phone, u.PasswordHash).Scan(&u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// Changed ? to $1
	query := `SELECT id, name, email, phone, password_hash, qr_code_url, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var u model.User
	var qrCode sql.NullString
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.PasswordHash, &qrCode, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if qrCode.Valid {
		u.QRCodeURL = qrCode.String
	}
	return &u, nil
}

func (r *SQLUserRepository) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	// Changed ? to $1
	query := `SELECT id, name, email, phone, password_hash, qr_code_url, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var u model.User
	var qrCode sql.NullString
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.PasswordHash, &qrCode, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if qrCode.Valid {
		u.QRCodeURL = qrCode.String
	}
	return &u, nil
}

func (r *SQLUserRepository) UpdateQRCode(ctx context.Context, userID uint64, qrURL string) error {
	// Changed ? to $1 and $2
	query := `UPDATE users SET qr_code_url = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, qrURL, userID)
	return err
}
