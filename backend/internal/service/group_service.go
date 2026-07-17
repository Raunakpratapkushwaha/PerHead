package service

import (
	"context"
	"errors"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/model"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/repository"
)

type GroupService interface {
	CreateGroup(ctx context.Context, req *model.CreateGroupRequest, creatorID int64) (*model.Group, error)
	AddMember(ctx context.Context, groupID, userID, requesterID int64) error
	RemoveMember(ctx context.Context, groupID, userID, requesterID int64) error
	GetGroup(ctx context.Context, id int64, requesterID int64) (*model.Group, error)
	ListGroups(ctx context.Context, userID int64) ([]model.Group, error)
}

type groupService struct {
	groupRepo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) GroupService {
	return &groupService{groupRepo: repo}
}

func (s *groupService) CreateGroup(ctx context.Context, req *model.CreateGroupRequest, creatorID int64) (*model.Group, error) {
	group := &model.Group{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.groupRepo.Create(ctx, group, creatorID); err != nil {
		return nil, err
	}
	return group, nil
}

func (s *groupService) AddMember(ctx context.Context, groupID, userID, requesterID int64) error {
	// Verify requester belongs to group to grant invite permission
	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return err
	}
	authorized := false
	for _, id := range members {
		if id == requesterID {
			authorized = true
			break
		}
	}
	if !authorized {
		return errors.New("unauthorized: must belong to the group to add members")
	}

	return s.groupRepo.AddMember(ctx, groupID, userID)
}

func (s *groupService) RemoveMember(ctx context.Context, groupID, userID, requesterID int64) error {
	// Verify requester is part of group
	members, err := s.groupRepo.GetMembers(ctx, groupID)
	if err != nil {
		return err
	}
	authorized := false
	for _, id := range members {
		if id == requesterID {
			authorized = true
			break
		}
	}
	if !authorized {
		return errors.New("unauthorized: must belong to the group to modify members")
	}

	return s.groupRepo.RemoveMember(ctx, groupID, userID)
}

func (s *groupService) GetGroup(ctx context.Context, id int64, requesterID int64) (*model.Group, error) {
	members, err := s.groupRepo.GetMembers(ctx, id)
	if err != nil {
		return nil, err
	}
	authorized := false
	for _, mID := range members {
		if mID == requesterID {
			authorized = true
			break
		}
	}
	if !authorized {
		return nil, errors.New("unauthorized to view group information")
	}

	return s.groupRepo.GetByID(ctx, id)
}

func (s *groupService) ListGroups(ctx context.Context, userID int64) ([]model.Group, error) {
	return s.groupRepo.ListUserGroups(ctx, userID)
}