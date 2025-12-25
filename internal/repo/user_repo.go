package repo

import (
	"my-chat/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	FindByPhone(phone string) (*model.User, error)
	FindByUuid(uuid string) (*model.User, error)
	FindUsersByIDs(ids []string) (map[string]*model.User, error)
	UpdateUser(user *model.User) error
}
type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Model(&model.User{}).
		Where("uuid = ?", user.Uuid).
		Updates(user).Error
}

func (r *userRepository) FindUsersByIDs(ids []string) (map[string]*model.User, error) {
	var users []*model.User
	if len(ids) == 0 {
		return make(map[string]*model.User), nil
	}
	err := r.db.Where("id IN (?)", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*model.User)
	for _, user := range users {
		userMap[user.Uuid] = user
	}
	return userMap, nil
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
