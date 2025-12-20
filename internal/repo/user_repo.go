package repo

import (
	"my-chat/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	FindByPhone(phone string) (*model.User, error)
	FindByUuid(uuid string) (*model.User, error)
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
func (r *userRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("telephone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) FindByUuid(uuid string) (*model.User, error) {
	var user model.User
	err := r.db.Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
