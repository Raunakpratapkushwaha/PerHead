package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/repository"
)

type ExpenseService interface {
	CreateExpense(ctx context.Context, req *model.CreateExpenseRequest, creatorID int64) (*model.Expense, error)
	GetGroupExpenses(ctx context.Context, groupID, userID int64) ([]model.Expense, error)
	GetGroupBalances(ctx context.Context, groupID, userID int64) (map[int64]int64, error)
}

type expenseService struct {
	expenseRepo repository.ExpenseRepository
	groupRepo   repository.GroupRepository
}

func NewExpenseService(eRepo repository.ExpenseRepository, gRepo repository.GroupRepository) ExpenseService {
	return &expenseService{expenseRepo: eRepo, groupRepo: gRepo}
}

func (s *expenseService) CreateExpense(ctx context.Context, req *model.CreateExpenseRequest, creatorID int64) (*model.Expense, error) {
	// 1. Verify the group exists and the payer & creator are members
	members, err := s.groupRepo.GetMembers(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, errors.New("group not found or has no members")
	}

	memberMap := make(map[int64]bool)
	for _, id := range members {
		memberMap[id] = true
	}

	if !memberMap[req.PayerID] {
		return nil, errors.New("payer is not a member of this group")
	}
	if !memberMap[creatorID] {
		return nil, errors.New("creator is not a member of this group")
	}

	// Validate that all split participants are part of this group
	for _, split := range req.Splits {
		if !memberMap[split.UserID] {
			return nil, fmt.Errorf("user %d is not a member of this group", split.UserID)
		}
	}

	// 2. Perform Splitting Engine Math
	splits, err := CalculateSplits(req.Amount, req.SplitType, req.Splits)
	if err != nil {
		return nil, fmt.Errorf("splitting engine failure: %w", err)
	}

	// 3. Assemble and persist the expense transactional payload
	expense := &model.Expense{
		GroupID:     req.GroupID,
		PayerID:     req.PayerID,
		Amount:      req.Amount,
		Description: req.Description,
		Category:    req.Category,
		SplitType:   req.SplitType,
		CreatedBy:   creatorID,
	}

	if err := s.expenseRepo.CreateWithSplits(ctx, expense, splits); err != nil {
		return nil, err
	}

	return expense, nil
}

// CalculateSplits converts user inputs into absolute splits (in cents) with penny-rounding safety.
func CalculateSplits(totalAmount int64, splitType model.SplitType, inputs []model.SplitInput) ([]model.ExpenseSplit, error) {
	if len(inputs) == 0 {
		return nil, errors.New("splits cannot be empty")
	}

	var splits []model.ExpenseSplit
	var allocatedSum int64

	switch splitType {
	case model.SplitEqual:
		// Divide total amount by participant count
		n := int64(len(inputs))
		baseAmount := totalAmount / n
		remainder := totalAmount % n // Discrepancy in cents

		for i, input := range inputs {
			finalAmount := baseAmount
			// Distribute remaining fractional cents 1 by 1
			if int64(i) < remainder {
				finalAmount++
			}
			splits = append(splits, model.ExpenseSplit{
				UserID: input.UserID,
				Amount: finalAmount,
			})
			allocatedSum += finalAmount
		}

	case model.SplitExact:
		// User states exact amounts. We must guarantee sum matches the parent total exactly.
		for _, input := range inputs {
			if input.Amount <= 0 {
				return nil, errors.New("exact split amount must be greater than 0")
			}
			splits = append(splits, model.ExpenseSplit{
				UserID: input.UserID,
				Amount: input.Amount,
			})
			allocatedSum += input.Amount
		}
		if allocatedSum != totalAmount {
			return nil, fmt.Errorf("sum of exact splits (%d cents) does not match total amount (%d cents)", allocatedSum, totalAmount)
		}

	case model.SplitPercentage:
		// Splitting based on percentages. We enforce that total percentage sums to 100%.
		var pctSum float64
		for _, input := range inputs {
			if input.Percentage < 0 {
				return nil, errors.New("percentages cannot be negative")
			}
			pctSum += input.Percentage
		}
		// Allowing a tiny precision margin due to floating point sum variance
		if pctSum < 99.99 || pctSum > 100.01 {
			return nil, fmt.Errorf("total percentage must equal 100%%, got %.2f%%", pctSum)
		}

		for _, input := range inputs {
			// Convert percent to scale of 0-10000 (basis points)
			bps := int64(input.Percentage * 100)
			userAmount := (totalAmount * bps) / 10000
			splits = append(splits, model.ExpenseSplit{
				UserID:     input.UserID,
				Amount:     userAmount,
				Percentage: input.Percentage,
			})
			allocatedSum += userAmount
		}

		// Adjust any rounding truncation to the first participant's split
		diff := totalAmount - allocatedSum
		if diff != 0 && len(splits) > 0 {
			splits[0].Amount += diff
		}

	case model.SplitShares:
		// Splitting by relative ratio shares (e.g., Alice: 1 share, Bob: 2 shares)
		var totalShares int
		for _, input := range inputs {
			if input.Share <= 0 {
				return nil, errors.New("shares must be greater than 0")
			}
			totalShares += input.Share
		}

		for _, input := range inputs {
			userAmount := (totalAmount * int64(input.Share)) / int64(totalShares)
			splits = append(splits, model.ExpenseSplit{
				UserID: input.UserID,
				Amount: userAmount,
				Share:  input.Share,
			})
			allocatedSum += userAmount
		}

		// Adjust any truncation to the first participant's split
		diff := totalAmount - allocatedSum
		if diff != 0 && len(splits) > 0 {
			splits[0].Amount += diff
		}

	default:
		return nil, fmt.Errorf("invalid split type: %s", splitType)
	}

	return splits, nil
}

func (s *expenseService) GetGroupExpenses(ctx context.Context, groupID, userID int64) ([]model.Expense, error) {
	// Verify user belongs to the group before revealing ledger
	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	authorized := false
	for _, id := range members {
		if id == userID {
			authorized = true
			break
		}
	}
	if !authorized {
		return nil, errors.New("unauthorized view request for this group")
	}

	return s.expenseRepo.GetByGroupID(ctx, groupID)
}

func (s *expenseService) GetGroupBalances(ctx context.Context, groupID, userID int64) (map[int64]int64, error) {
	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	authorized := false
	for _, id := range members {
		if id == userID {
			authorized = true
			break
		}
	}
	if !authorized {
		return nil, errors.New("unauthorized balance request for this group")
	}

	return s.expenseRepo.GetBalances(ctx, groupID)
}