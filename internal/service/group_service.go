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
	userRepo  repo.UserRepository
}

func NewGroupService(groupRepo repo.GroupRepository, userRepo repo.UserRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
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

type GroupMemberResp struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Role     int    `json:"role"`
	JoinedAt string `json:"joined_at"`
}

func (s *GroupService) GetGroupInfo(groupId string) (*model.Group, error) {
	return s.groupRepo.FindGroup(groupId)
}
func (s *GroupService) GetGroupMembers(groupId string) ([]GroupMemberResp, error) {
	members, err := s.groupRepo.GetGroupMembers(groupId)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return []GroupMemberResp{}, nil
	}
	var userIds []string
	for _, member := range members {
		userIds = append(userIds, member.UserId)
	}
	userMap, err := s.userRepo.FindUsersByIDs(userIds)
	if err != nil {
		return nil, err
	}
	var resp []GroupMemberResp
	for _, member := range members {
		user, exists := userMap[member.UserId]
		displayNickname := ""
		avatar := ""
		if exists {
			avatar = user.Avatar
			if member.Nickname != "" {
				displayNickname = member.Nickname
			} else {
				displayNickname = user.Nickname
			}
		} else {
			displayNickname = "未知用户"
		}
		resp = append(resp, GroupMemberResp{
			UserId:   member.UserId,
			Nickname: displayNickname,
			Avatar:   avatar,
			Role:     member.Role,
			JoinedAt: member.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return resp, nil
}
func (s *GroupService) LoadMyJoinedGroup(userId string) (*model.Group, error) {
	return s.groupRepo.FindGroup(userId)
}
func (s *GroupService) LeaveGroup(groupId, userId string) error {
	group, err := s.groupRepo.FindGroup(groupId)
	if err != nil {
		return err
	}
	if group.OwnerId == userId {
		return errno.New(30005, "群主不能直接退群，需要转让或解散群")
	}
	return s.groupRepo.RemoveMember(groupId, userId)
}
func (s *GroupService) KickMember(operatorId, groupId, userId string) error {
	group, err := s.groupRepo.FindGroup(groupId)
	if err != nil {
		return err
	}
	if group.OwnerId != operatorId {
		return errno.New(30006, "权限不足，只有群主可以移除")
	}
	return s.groupRepo.RemoveMember(groupId, userId)
}
func (s *GroupService) DismissGroup(operatorId, groupId string) error {
	group, err := s.groupRepo.FindGroup(groupId)
	if err != nil {
		return err
	}
	if group.OwnerId != operatorId {
		return errno.New(30006, "权限不足，只有群主可以移除")
	}
	return s.groupRepo.DeleteGroup(groupId)
}
