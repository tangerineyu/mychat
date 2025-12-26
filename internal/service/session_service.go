package service

import (
	//"my-chat/internal/model"
	"my-chat/internal/repo"
)

type SessionService struct {
	sessionRepo repo.SessionRepository
	groupRepo   repo.GroupRepository
	userRepo    repo.UserRepository
}

func NewSessionService(sessionRepo repo.SessionRepository, groupRepo repo.GroupRepository, userRepo repo.UserRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		groupRepo:   groupRepo,
		userRepo:    userRepo,
	}
}

type SessionDto struct {
	TargetId  string `json:"target_id"`
	Type      int    `json:"type"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	LastMsg   string `json:"last_msg"`
	LastTime  int64  `json:"last_time"`
	UnreadCnt int    `json:"unread_cnt"`
}

func (s *SessionService) GetUserSessions(userId string) ([]SessionDto, error) {
	sessions, err := s.sessionRepo.GetList(userId)
	if err != nil {
		return nil, err
	}
	var userIds []string
	var groupIds []string
	for _, session := range sessions {
		if session.Type == 1 {
			userIds = append(userIds, session.TargetId)
		} else {
			groupIds = append(groupIds, session.TargetId)
		}
	}
	var result []SessionDto
	for _, sess := range sessions {
		name := "未知"
		avatar := ""
		if sess.Type == 1 {
			u, _ := s.userRepo.FindByUuid(sess.TargetId)
			if u != nil {
				name = u.Nickname
				avatar = u.Avatar
			}
		} else {
			g, _ := s.groupRepo.FindGroup(sess.TargetId)
			if g != nil {
				name = g.Name
				avatar = g.Avatar
			}
		}
		result = append(result, SessionDto{
			TargetId:  sess.TargetId,
			Type:      sess.Type,
			Name:      name,
			Avatar:    avatar,
			LastMsg:   sess.LastMsg,
			LastTime:  sess.LastTime,
			UnreadCnt: sess.UnreadCnt,
		})
	}
	return result, nil
}
