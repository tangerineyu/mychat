package service

import (
	"errors"
	"my-chat/internal/model"
	"my-chat/internal/repo"
	"my-chat/pkg/errno"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupService struct {
	groupRepo repo.GroupRepository
}

func NewGroupService(groupRepo repo.GroupRepository) *GroupService {
	return &GroupService{groupRepo: groupRepo}
}
func (s *GroupService) CreateGroup(ownerId, name string) (*model.Group, error) {
	groupId := "G" + uuid.New().String()
	newGroup := &model.Group{
		Uuid:    groupId,
		Name:    name,
		OwnerId: ownerId,
		Notice:  "欢迎加入群聊",
	}
	ownerMember := &model.GroupMember{
		GroupId:  groupId,
		UserId:   ownerId,
		Role:     model.RoleOwner,
		Nickname: "群主",
	}
	if err := s.groupRepo.CreateGroup(newGroup, ownerMember); err != nil {
		return nil, err
	}
	return newGroup, nil
}
func (s *GroupService) JoinGroup(groupId, userId string) error {
	_, err := s.groupRepo.FindGroup(groupId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrGroupNotFound
		}
		return err
	}
	isMember, err := s.groupRepo.IsMember(groupId, userId)
	if err != nil {
		return err
	}
	if isMember {
		return errno.New(30404, "已经是群成员，不能重复加入")
	}
	newMember := &model.GroupMember{
		GroupId: groupId,
		UserId:  userId,
		Role:    model.RoleMember,
	}
	return s.groupRepo.AddMember(newMember)
}
