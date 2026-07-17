package service

import (
	"context"
	"sort"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
)

type SettlementService interface {
	GetSimplifiedDebts(ctx context.Context, groupID, userID int64) ([]model.Settlement, error)
}

type settlementService struct {
	expenseService ExpenseService
}

func NewSettlementService(es ExpenseService) SettlementService {
	return &settlementService{expenseService: es}
}

// userBalance is a helper struct for the simplification algorithm
type userBalance struct {
	userID int64
	amount int64 // Always stored as a positive absolute value here
}

func (s *settlementService) GetSimplifiedDebts(ctx context.Context, groupID, userID int64) ([]model.Settlement, error) {
	// 1. Fetch the raw net balances from our ExpenseService
	balances, err := s.expenseService.GetGroupBalances(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}

	// 2. Separate into Debtors and Creditors
	var debtors []userBalance
	var creditors []userBalance

	for uid, balance := range balances {
		if balance < 0 {
			// Convert to positive absolute value for easier math
			debtors = append(debtors, userBalance{userID: uid, amount: -balance})
		} else if balance > 0 {
			creditors = append(creditors, userBalance{userID: uid, amount: balance})
		}
	}

	// 3. Sort descending to optimize the greedy match (largest debts match with largest credits)
	sort.Slice(debtors, func(i, j int) bool { return debtors[i].amount > debtors[j].amount })
	sort.Slice(creditors, func(i, j int) bool { return creditors[i].amount > creditors[j].amount })

	var settlements []model.Settlement
	i, j := 0, 0

	// 4. Greedily resolve debts
	for i < len(debtors) && j < len(creditors) {
		debtor := &debtors[i]
		creditor := &creditors[j]

		// The settlement amount is the minimum of what the debtor owes and what the creditor is owed
		settleAmount := debtor.amount
		if creditor.amount < debtor.amount {
			settleAmount = creditor.amount
		}

		// Record the transaction
		settlements = append(settlements, model.Settlement{
			FromUserID: debtor.userID,
			ToUserID:   creditor.userID,
			Amount:     settleAmount,
		})

		// Deduct the settled amount from both parties
		debtor.amount -= settleAmount
		creditor.amount -= settleAmount

		// Move the pointers if a person's balance is fully resolved
		if debtor.amount == 0 {
			i++
		}
		if creditor.amount == 0 {
			j++
		}
	}

	return settlements, nil
}
