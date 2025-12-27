package service

import (
	//"my-chat/internal/model"
	"my-chat/internal/repo"
)

type AdminService struct {
	adminRepo repo.AdminRepository
}

func NewAdminService(adminRepo repo.AdminRepository) *AdminService {
	return &AdminService{adminRepo: adminRepo}
}
func (s *AdminService) GetUserList(page, limit int) (map[string]interface{}, error) {
	users, total, err := s.adminRepo.GetAllUsers(page, limit)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"list":  users,
		"total": total,
	}, nil
}
func (s *AdminService) BanUser(uuid string) error {
	return s.adminRepo.UpdateUserStatus(uuid, 2)
}
func (s *AdminService) UnBanUser(uuid string) error {
	return s.adminRepo.UpdateUserStatus(uuid, 1)
}
func (s *AdminService) GetGroupList(page, limit int) (map[string]interface{}, error) {
	groups, total, err := s.adminRepo.GetAllGroups(page, limit)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"list":  groups,
		"total": total,
	}, nil
}
func (s *AdminService) BanGroup(uuid string) error {
	return s.adminRepo.UpdateUserStatus(uuid, 2)
}
func (s *AdminService) UnBanGroup(uuid string) error {
	return s.adminRepo.UpdateUserStatus(uuid, 1)
}
