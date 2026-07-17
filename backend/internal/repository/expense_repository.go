package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
)

type ExpenseRepository interface {
	CreateWithSplits(ctx context.Context, expense *model.Expense, splits []model.ExpenseSplit) error
	GetByGroupID(ctx context.Context, groupID int64) ([]model.Expense, error)
	GetSplitsByExpenseID(ctx context.Context, expenseID int64) ([]model.ExpenseSplit, error)
	GetBalances(ctx context.Context, groupID int64) (map[int64]int64, error)
}

type expenseRepo struct {
	db *sqlx.DB
}

func NewExpenseRepository(db *sqlx.DB) ExpenseRepository {
	return &expenseRepo{db: db}
}

func (r *expenseRepo) CreateWithSplits(ctx context.Context, expense *model.Expense, splits []model.ExpenseSplit) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert the parent expense
	query := `
		INSERT INTO expenses (group_id, payer_id, amount, description, category, split_type, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`
	err = tx.QueryRowContext(ctx, query, expense.GroupID, expense.PayerID, expense.Amount, expense.Description, expense.Category, expense.SplitType, expense.CreatedBy).Scan(&expense.ID, &expense.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert expense: %w", err)
	}

	// 2. Bulk insert the calculated splits
	splitQuery := `
		INSERT INTO expense_splits (expense_id, user_id, amount, percentage, share)
		VALUES ($1, $2, $3, $4, $5)`
	for i := range splits {
		splits[i].ExpenseID = expense.ID
		_, err = tx.ExecContext(ctx, splitQuery, splits[i].ExpenseID, splits[i].UserID, splits[i].Amount, splits[i].Percentage, splits[i].Share)
		if err != nil {
			return fmt.Errorf("failed to insert split for user %d: %w", splits[i].UserID, err)
		}
	}

	return tx.Commit()
}

func (r *expenseRepo) GetByGroupID(ctx context.Context, groupID int64) ([]model.Expense, error) {
	var expenses []model.Expense
	query := `SELECT * FROM expenses WHERE group_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &expenses, query, groupID)
	return expenses, err
}

func (r *expenseRepo) GetSplitsByExpenseID(ctx context.Context, expenseID int64) ([]model.ExpenseSplit, error) {
	var splits []model.ExpenseSplit
	query := `SELECT * FROM expense_splits WHERE expense_id = $1`
	err := r.db.SelectContext(ctx, &splits, query, expenseID)
	return splits, err
}

// GetBalances aggregates who owes what in a group.
// Returns a map where key = user_id and value = balance in cents.
// Balance = (Total Paid by User) - (Total Owed by User)
func (r *expenseRepo) GetBalances(ctx context.Context, groupID int64) (map[int64]int64, error) {
	balances := make(map[int64]int64)

	type userBalance struct {
		UserID  int64 `db:"user_id"`
		Balance int64 `db:"balance"`
	}

	var results []userBalance
	query := `
		WITH expense_paid AS (
			SELECT payer_id AS user_id, COALESCE(SUM(amount), 0) AS total_paid
			FROM expenses
			WHERE group_id = $1
			GROUP BY payer_id
		),
		expense_owed AS (
			SELECT es.user_id, COALESCE(SUM(es.amount), 0) AS total_owed
			FROM expense_splits es
			JOIN expenses e ON es.expense_id = e.id
			WHERE e.group_id = $1
			GROUP BY es.user_id
		),
		payments_sent AS (
			SELECT payer_id AS user_id, COALESCE(SUM(amount), 0) AS total_sent
			FROM payments
			WHERE group_id = $1
			GROUP BY payer_id
		),
		payments_received AS (
			SELECT payee_id AS user_id, COALESCE(SUM(amount), 0) AS total_received
			FROM payments
			WHERE group_id = $1
			GROUP BY payee_id
		),
		group_members AS (
			SELECT user_id FROM group_members WHERE group_id = $1
		)
		SELECT 
			gm.user_id,
			(COALESCE(ep.total_paid, 0) + COALESCE(ps.total_sent, 0)) - 
			(COALESCE(eo.total_owed, 0) + COALESCE(pr.total_received, 0)) AS balance
		FROM group_members gm
		LEFT JOIN expense_paid ep ON gm.user_id = ep.user_id
		LEFT JOIN expense_owed eo ON gm.user_id = eo.user_id
		LEFT JOIN payments_sent ps ON gm.user_id = ps.user_id
		LEFT JOIN payments_received pr ON gm.user_id = pr.user_id;
	`

	if err := r.db.SelectContext(ctx, &results, query, groupID); err != nil {
		return nil, err
	}

	for _, r := range results {
		balances[r.UserID] = r.Balance
	}

	return balances, nil
}