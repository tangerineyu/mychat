package dao

import (
	"fmt"
	"my-chat/internal/config"
	"my-chat/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//var DB *gorm.DB

func NewMySQL(c *config.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&model.User{},
		&model.Group{},
		&model.GroupMember{},
		&model.Message{},
		&model.Contact{},
		&model.ContactApply{},
		&model.Session{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
