package repo

import (
	"my-chat/internal/model"

	"gorm.io/gorm"
)

type AdminRepository interface {
	GetAllUsers(page, limit int) ([]*model.User, int64, error)
	UpdateUserStatus(uuid string, status int) error
	GetAllGroups(page, limit int) ([]*model.Group, int64, error)
	UpdateGroupStatus(uuid string, status int) error
}
type adminRepository struct {
	db *gorm.DB
}

func (a *adminRepository) GetAllGroups(page, limit int) ([]*model.Group, int64, error) {
	var groups []*model.Group
	var total int64
	offset := (page - 1) * limit
	a.db.Model(&model.Group{}).Count(&total)
	err := a.db.Offset(offset).Limit(limit).Find(&groups).Error
	return groups, total, err
}

func (a *adminRepository) UpdateGroupStatus(uuid string, status int) error {
	return a.db.Model(&model.Group{}).Where("uuid = ?", uuid).Update("status", status).Error
}

func (a *adminRepository) UpdateUserStatus(uuid string, status int) error {
	return a.db.Model(&model.User{}).Where("uuid = ?", uuid).Update("status", status).Error
}

func (a *adminRepository) GetAllUsers(page, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64
	offset := (page - 1) * limit
	a.db.Model(&model.User{}).Count(&total)
	err := a.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}
