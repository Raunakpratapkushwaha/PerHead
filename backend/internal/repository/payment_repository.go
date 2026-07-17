package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
)

type PaymentRepository interface {
	Record(ctx context.Context, p *model.Payment) error
	GetByGroupID(ctx context.Context, groupID int64) ([]model.Payment, error)
}

type paymentRepo struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Record(ctx context.Context, p *model.Payment) error {
	query := `
		INSERT INTO payments (group_id, payer_id, payee_id, amount, notes, payment_method, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, p.GroupID, p.PayerID, p.PayeeID, p.Amount, p.Notes, p.PaymentMethod, p.CreatedBy).
		Scan(&p.ID, &p.CreatedAt)
}

func (r *paymentRepo) GetByGroupID(ctx context.Context, groupID int64) ([]model.Payment, error) {
	var payments []model.Payment
	query := `SELECT * FROM payments WHERE group_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &payments, query, groupID)
	return payments, err
}