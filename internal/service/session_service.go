package service

import (
	//"my-chat/internal/model"
	"my-chat/internal/model"
	"my-chat/internal/repo"
	"my-chat/pkg/zlog"
	"sort"

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
	//获取私聊会话
	//可能返回群会话，在下面遍历的时候过滤掉
	sessions, err := s.sessionRepo.GetList(userId)
	if err != nil {
		return nil, err
	}
	//获取群聊会话
	groupList, err := s.groupRepo.GetUserJoinedGroups(userId)
	if err != nil {
		return nil, err
	}
	var friendIds []string
	//var groupIds []string
	//移除了groupIds的收集，直接调用了GetUserJoinedGroups，这个方法返回的就是Group结构体
	//里面已经包含了Name， Avatar，LAstMsg，LastTime
	for _, session := range sessions {
		if session.Type == 1 {
			friendIds = append(friendIds, session.TargetId)
		}
	}
	userMap, err := s.userRepo.FindUsersByIDs(friendIds)
	//移除了groupMap，直接遍历groupList转换为SessionDto就可以
	if err != nil {
		return nil, err
	}
	var result []SessionDto
	for _, sess := range sessions {
		if sess.Type == 1 {
			name := "未知"
			avatar := ""
			if user, ok := userMap[sess.TargetId]; ok {
				name = user.Nickname
				avatar = user.Avatar
			}
			result = append(result, SessionDto{
				TargetId:  sess.TargetId,
				Type:      1,
				Name:      name,
				Avatar:    avatar,
				LastMsg:   sess.LastMsg,
				LastTime:  sess.LastTime,
				UnreadCnt: sess.UnreadCnt,
			})
		}
	}
	//处理群聊
	for _, group := range groupList {
		result = append(result, SessionDto{
			TargetId:  group.Uuid,
			Type:      2,
			Name:      group.Name,
			Avatar:    group.Avatar,
			LastMsg:   group.LastMsg,
			LastTime:  group.LastTime,
			UnreadCnt: 0,
		})
	}
	//排序，将私聊与群聊混合在一起，按照时间最新的排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastTime > result[j].LastTime
	})
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
