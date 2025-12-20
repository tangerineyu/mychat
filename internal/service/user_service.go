package service

import (
	"errors"
	"my-chat/internal/model"
	"my-chat/internal/repo"
	"my-chat/pkg/errno"
	"my-chat/pkg/util/password"
	"my-chat/pkg/util/token"
	"my-chat/pkg/zlog"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
func (s *UserService) Register(phone, rawPassword, nickName string) error {
	_, err := s.userRepo.FindByPhone(phone)
	if err == nil {
		zlog.Error("此手机号已经注册",
			zap.String("phone", phone),
			zap.Error(err))
		return errno.ErrUserAlreadyExist
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hashPwd, err := password.HashPassword(rawPassword)
	if err != nil {
		return err
	}
	newUser := &model.User{
		Uuid:      "U" + uuid.New().String(),
		Telephone: phone,
		Password:  hashPwd,
		Nickname:  nickName,
		Status:    1,
		Avatar:    "default_avatar.png",
	}
	return s.userRepo.CreateUser(newUser)
}
func (s *UserService) Login(phone, rawPassword string) (string, string, *model.User, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.Error("用户不存在",
				zap.String("phone", phone),
				zap.Error(err))
			return "", "", nil, errno.ErrUserBanned
		}
		return "", "", nil, err
	}
	if !password.CheckPassword(rawPassword, user.Password) {
		return "", "", nil, errors.New("密码错误")
	}
	accessToken, refreshToken, err := token.GenerateTokenPair(user.Uuid)
	if err != nil {
		return "", "", nil, errors.New("token生成失败")
	}
	return accessToken, refreshToken, user, nil
}
func (s *UserService) RefreshToken(refreshTokenStr string) (string, string, error) {
	claims, err := token.ParseRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", errors.New("Refresh Token 无效或已过期，请重新登录")
	}
	user, err := s.userRepo.FindByUuid(claims.UserId)
	if err != nil {
		return "", "", errors.New("用户不存在")
	}
	if user.Status != 1 {
		return "", "", errors.New("用户已被禁用")
	}
	newAccess, newRefresh, err := token.GenerateTokenPair(claims.UserId)
	if err != nil {
		return "", "", errors.New("token生成失败")
	}
	return newAccess, newRefresh, nil
}
