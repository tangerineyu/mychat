package service

import (
	//"my-chat/internal/model"
	"my-chat/internal/model"
	"my-chat/internal/repo"
	"my-chat/pkg/zlog"

	"go.uber.org/zap"
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
	//先查redis
	cacheList, err := s.sessionRepo.GetListFromCache(userId)
	if err == nil && len(cacheList) > 0 {
		//缓存命中，直接组装返回，不用mysql
		var result []SessionDto
		for _, v := range cacheList {
			result = append(result, SessionDto{
				TargetId:  v.TargetId,
				Type:      v.Type,
				Name:      v.Name,
				Avatar:    v.Avatar,
				LastMsg:   v.LastMsg,
				LastTime:  v.LastTime,
				UnreadCnt: v.UnreadCnt,
			})
		}
		zlog.Info("Session list hit cache",
			zap.String("userId", userId))
		return result, nil
	}
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
	userMap, err := s.userRepo.FindUsersByIDs(userIds)
	if err != nil {
		return nil, err
	}
	groupMap, err := s.groupRepo.FindGroupsByIds(groupIds)
	if err != nil {
		return nil, err
	}
	var result []SessionDto
	for _, sess := range sessions {
		name := "未知"
		avatar := ""
		if sess.Type == 1 {
			if user, ok := userMap[sess.TargetId]; ok {
				name = user.Nickname
				avatar = user.Avatar
			}
		} else {
			if group, ok := groupMap[sess.TargetId]; ok {
				name = group.Name
				avatar = group.Avatar
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
	//查完mysql，异步填回redis
	go func() {
		for _, dto := range result {
			cacheItem := &model.SessionCache{
				TargetId:  dto.TargetId,
				Type:      dto.Type,
				Name:      dto.Name,
				Avatar:    dto.Avatar,
				LastMsg:   dto.LastMsg,
				LastTime:  dto.LastTime,
				UnreadCnt: dto.UnreadCnt,
			}
			_ = s.sessionRepo.SaveSessionToCache(userId, cacheItem)
		}
	}()
	return result, nil
}
func (s *SessionService) UpsertSession(session *model.Session) error {

	return s.sessionRepo.UpsertSession(session)
}
