package service

import (
	"context"
	"errors"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/repository"
)

type PaymentService interface {
	RecordPayment(ctx context.Context, groupID, creatorID int64, req *model.RecordPaymentRequest) (*model.Payment, error)
	GetGroupPayments(ctx context.Context, groupID, userID int64) ([]model.Payment, error)
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	groupRepo   repository.GroupRepository
}

func NewPaymentService(pRepo repository.PaymentRepository, gRepo repository.GroupRepository) PaymentService {
	return &paymentService{paymentRepo: pRepo, groupRepo: gRepo}
}

func (s *paymentService) RecordPayment(ctx context.Context, groupID, creatorID int64, req *model.RecordPaymentRequest) (*model.Payment, error) {
	if req.PayerID == req.PayeeID {
		return nil, errors.New("cannot pay yourself back")
	}

	// Verify requester is a member of the group
	isCreatorMember, err := s.groupRepo.IsMember(ctx, groupID, creatorID)
	if err != nil || !isCreatorMember {
		return nil, errors.New("unauthorized: you do not belong to this group")
	}

	// Verify both parties involved in the payment belong to the group
	isPayerMember, err := s.groupRepo.IsMember(ctx, groupID, req.PayerID)
	if err != nil || !isPayerMember {
		return nil, errors.New("payer is not a member of this group")
	}

	isPayeeMember, err := s.groupRepo.IsMember(ctx, groupID, req.PayeeID)
	if err != nil || !isPayeeMember {
		return nil, errors.New("payee is not a member of this group")
	}

	payment := &model.Payment{
		GroupID:       groupID,
		PayerID:       req.PayerID,
		PayeeID:       req.PayeeID,
		Amount:        req.Amount,
		Notes:         req.Notes,
		PaymentMethod: req.PaymentMethod,
		CreatedBy:     creatorID,
	}

	if err := s.paymentRepo.Record(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentService) GetGroupPayments(ctx context.Context, groupID, userID int64) ([]model.Payment, error) {
	isMember, err := s.groupRepo.IsMember(ctx, groupID, userID)
	if err != nil || !isMember {
		return nil, errors.New("unauthorized view request for this group")
	}

	return s.paymentRepo.GetByGroupID(ctx, groupID)
}