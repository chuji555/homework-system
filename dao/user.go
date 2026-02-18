package dao

import (
	"github.com/chuji555/homework-system/models"
	"gorm.io/gorm"
)

// 注册用户
func CreateUser(user *models.User) error {
	return DB.Create(user).Error
}

// 根据用户名查询用户（登录用）
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// 根据ID查询用户（刷新Token用）
func GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	err := DB.Where("id = ?", userID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// 软删除用户（注销账号）
func DeleteUser(userID int64) error {
	return DB.Delete(&models.User{}, userID).Error
}
