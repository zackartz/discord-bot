package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"synergy/models"
	"time"
)

var (
	db *gorm.DB
)

func Init() {
	var err error

	DBURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	db, err = gorm.Open(mysql.Open(DBURI), &gorm.Config{})
	if err != nil {
		log.Panicf("could not connect to DB %v", err)
	}

	db.Debug().AutoMigrate(&models.RoleMessage{}, &models.Role{})
}

func CreateRoleMessage(id int64, channelID int64, guildID int64) error {
	var roleMsg models.RoleMessage
	err := db.Model(&models.RoleMessage{}).Where("channel_id = ?", channelID).Take(&roleMsg).Error
	if err != gorm.ErrRecordNotFound {
		return errors.New("only one message per channel")
	}

	msg := &models.RoleMessage{
		ID:        id,
		ChannelID: channelID,
		GuildID:   guildID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = db.Model(&msg).Create(&msg).Error
	return err
}

func GetRoleMessageByChannelID(id int64) (*models.RoleMessage, error) {
	var roleMsg models.RoleMessage
	err := db.Model(&models.RoleMessage{}).Where("channel_id = ?", id).Take(&roleMsg).Error
	if err != nil {
		return nil, err
	}
	roleMsg.Roles, err = GetRolesForMessage(roleMsg.ID)
	if err != nil {
		return nil, err
	}
	return &roleMsg, err
}

func GetRoleMessageByMessageID(id int64) (*models.RoleMessage, error) {
	var roleMsg models.RoleMessage
	err := db.Model(&models.RoleMessage{}).Where("id = ?", id).Take(&roleMsg).Error
	if err != nil {
		return nil, err
	}
	roleMsg.Roles, err = GetRolesForMessage(roleMsg.ID)
	if err != nil {
		return nil, err
	}
	return &roleMsg, err
}

func GetRolesForMessage(id int64) ([]*models.Role, error) {
	var roles []*models.Role
	err := db.Model(&models.Role{}).Where("message_id = ?", id).Find(&roles).Error
	return roles, err
}

func AddRoleToMessage(roleID, messageID, emojiID int64, emoji, name string) error {
	role := &models.Role{
		ID:        roleID,
		Name:      name,
		Emoji:     emoji,
		EmojiID:   emojiID,
		MessageID: messageID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return db.Model(&models.Role{}).Create(&role).Error
}

func RemoveRoleFromMessge(roleID int64) error {
	return db.Model(&models.Role{}).Where("id = ?", roleID).Delete(&models.Role{}).Error
}

func GetRoleByID(roleID int64) (*models.Role, error) {
	var role *models.Role
	return role, db.Model(&models.Role{}).Where("id = ?", roleID).Take(&role).Error
}

func GetRoleByEmoji(emoji string) (*models.Role, error) {
	role := &models.Role{}
	err := db.Model(&models.Role{}).Where("emoji = ?", emoji).Take(&role).Error
	if err != nil {
		return role, err
	}
	return role, err
}
