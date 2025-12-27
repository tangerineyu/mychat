package service

import (
	"errors"
	"my-chat/internal/model"
	"my-chat/internal/repo"
	"my-chat/pkg/zlog"

	"go.uber.org/zap"
)

type ContactService struct {
	contactRepo repo.ContactRepository
	userRepo    repo.UserRepository
}

func NewContactService(contactRepo repo.ContactRepository, userRepo repo.UserRepository) *ContactService {
	return &ContactService{contactRepo: contactRepo, userRepo: userRepo}
}
func (s *ContactService) AddFriendApply(userId, targetId, msg string) error {
	if userId == targetId {
		zlog.Warn("不能添加自己为好友",
			zap.String("userId", userId))
		return errors.New("不能添加自己为好友")
	}
	apply := &model.ContactApply{
		UserId:   userId,
		TargetId: targetId,
		Msg:      msg,
		Status:   0,
	}
	return s.contactRepo.CreateApply(apply)
}
func (s *ContactService) RefuseApply(applyId uint) error {
	return s.contactRepo.UpdateApplyStatus(applyId, "2")
}
func (s *ContactService) RemoveFriend(userId, targetId string) error {
	return s.contactRepo.DeleteFriend(userId, targetId)
}
func (s *ContactService) GetApplyList(userId string) ([]*model.ContactApply, error) {
	return s.contactRepo.GetApplyList(userId)
}
func (s *ContactService) AgreeFriend(applyId uint) error {
	apply, err := s.contactRepo.FindApply(applyId)
	if err != nil {
		return err
	}
	if apply.Status != 0 {
		return errors.New("好友申请已处理，不能重复操作")
	}
	contactA := &model.Contact{
		OwnerId:  apply.UserId,
		TargetId: apply.TargetId,
		Type:     1,
	}
	contactB := &model.Contact{
		OwnerId:  apply.TargetId,
		TargetId: apply.UserId,
		Type:     1,
	}
	return s.contactRepo.AddFriend(apply.ID, contactA, contactB)
}
func (s *ContactService) GetContactList(userId string) ([]map[string]interface{}, error) {
	contacts, err := s.contactRepo.GetContacts(userId)
	if err != nil {
		return nil, err
	}
	var friendIds []string
	for _, c := range contacts {
		friendIds = append(friendIds, c.TargetId)
	}
	userMap, err := s.userRepo.FindUsersByIDs(friendIds)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for _, c := range contacts {
		user := userMap[c.TargetId]
		if user != nil {
			result = append(result, map[string]interface{}{
				"user_id":  user.Uuid,
				"nickname": user.Nickname,
				"avatar":   user.Avatar,
				"desc":     c.Desc,
			})
		}
	}
	return result, nil
}
func (s *ContactService) BlackContact(userId, targetId string) error {
	return s.contactRepo.UpdateContactType(userId, targetId, model.ContactTypeBlack)
}
func (s *ContactService) UnBlackContact(userId, targetId string) error {
	return s.contactRepo.UpdateContactType(userId, targetId, model.ContactTypeFriend)
}
