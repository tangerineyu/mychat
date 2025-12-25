package repo

import (
	"my-chat/internal/model"

	"gorm.io/gorm"
)

type ContactRepository interface {
	CreateApply(apply *model.ContactApply) error
	GetApplyList(userId string) ([]*model.ContactApply, error)
	GetContacts(ownerId string) ([]*model.Contact, error)
	FindApply(id uint) (*model.ContactApply, error)
	AddFriend(applyId uint, c1, c2 *model.Contact) error
}

type contactRepository struct {
	db *gorm.DB
}

func (c *contactRepository) CreateApply(apply *model.ContactApply) error {
	return c.db.Create(apply).Error
}

func (c *contactRepository) GetApplyList(userId string) ([]*model.ContactApply, error) {
	var list []*model.ContactApply
	err := c.db.Where("target_id = ?", userId).Order("created_at desc").Find(&list).Error
	return list, err
}

func (c *contactRepository) GetContacts(ownerId string) ([]*model.Contact, error) {
	var list []*model.Contact
	err := c.db.Where("owner_id = ?", ownerId).Find(&list).Error
	return list, err
}

func (c *contactRepository) FindApply(id uint) (*model.ContactApply, error) {
	var apply model.ContactApply
	err := c.db.First(&apply, id).Error
	return &apply, err
}

func (c *contactRepository) AddFriend(applyId uint, c1, c2 *model.Contact) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.ContactApply{}).
			Where("id = ?", applyId).
			Update("status", 1).Error; err != nil {
			return err
		}
		if err := tx.Create(c1).Error; err != nil {
			return err
		}
		if err := tx.Create(c2).Error; err != nil {
			return err
		}
		return nil
	})
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}
