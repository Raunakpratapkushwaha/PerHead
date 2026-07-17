package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
)

type GroupRepository interface {
	Create(ctx context.Context, group *model.Group, creatorID int64) error
	AddMember(ctx context.Context, groupID, userID int64) error
	RemoveMember(ctx context.Context, groupID, userID int64) error
	GetByID(ctx context.Context, id int64) (*model.Group, error)
	GetMembers(ctx context.Context, groupID int64) ([]int64, error)
	ListUserGroups(ctx context.Context, userID int64) ([]model.Group, error)
	IsMember(ctx context.Context, groupID, userID int64) (bool, error) // <-- Add this line here!
}

type groupRepo struct {
	db *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) GroupRepository {
	return &groupRepo{db: db}
}

func (r *groupRepo) Create(ctx context.Context, group *model.Group, creatorID int64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO groups (name, description) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err = tx.QueryRowContext(ctx, query, group.Name, group.Description).Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert group: %w", err)
	}

	memberQuery := `INSERT INTO group_members (group_id, user_id) VALUES ($1, $2)`
	_, err = tx.ExecContext(ctx, memberQuery, group.ID, creatorID)
	if err != nil {
		return fmt.Errorf("failed to add creator as member: %w", err)
	}

	return tx.Commit()
}

func (r *groupRepo) AddMember(ctx context.Context, groupID, userID int64) error {
	query := `INSERT INTO group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	return err
}

func (r *groupRepo) RemoveMember(ctx context.Context, groupID, userID int64) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	return err
}

func (r *groupRepo) GetByID(ctx context.Context, id int64) (*model.Group, error) {
	var g model.Group
	query := `SELECT * FROM groups WHERE id = $1`
	err := r.db.GetContext(ctx, &g, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &g, err
}

func (r *groupRepo) GetMembers(ctx context.Context, groupID int64) ([]int64, error) {
	var userIDs []int64
	query := `SELECT user_id FROM group_members WHERE group_id = $1`
	err := r.db.SelectContext(ctx, &userIDs, query, groupID)
	return userIDs, err
}

func (r *groupRepo) ListUserGroups(ctx context.Context, userID int64) ([]model.Group, error) {
	var groups []model.Group
	query := `
		SELECT g.* FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1
		ORDER BY g.updated_at DESC`
	err := r.db.SelectContext(ctx, &groups, query, userID)
	return groups, err
}

func (r *groupRepo) IsMember(ctx context.Context, groupID, userID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)`
	err := r.db.GetContext(ctx, &exists, query, groupID, userID)
	return exists, err
}