package service

import (
	"github.com/chuji555/homework-system/dao"
	"github.com/chuji555/homework-system/models"
	"github.com/chuji555/homework-system/pkg/errcode"
	"github.com/chuji555/homework-system/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// 注册业务（密码加密+数据校验）
func Register(username, password, nickname, department string) (*models.User, errcode.ErrCode) {
	// 校验参数
	if username == "" || password == "" || nickname == "" || department == "" {
		return nil, errcode.ParamError
	}

	// 检查用户名是否已存在
	existingUser, _ := dao.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errcode.ParamError // 可自定义"用户名已存在"错误码，这里简化
	}

	// 密码加密（加盐哈希）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errcode.DBError
	}
	// 创建用户（默认角色是student，可手动改admin）
	user := &models.User{
		Username:   username,
		Password:   string(hashedPassword),
		Nickname:   nickname,
		Role:       models.Student,
		Department: models.Department(department),
	}
	if err := dao.CreateUser(user); err != nil {
		return nil, errcode.DBError
	}
	return user, errcode.Success
}

// 登录业务（密码校验+生成双Token）
func Login(username, password string) (accessToken, refreshToken string, user *models.User, errCode errcode.ErrCode) {
	// 1. 查询用户
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		return "", "", nil, errcode.DBError
	}
	if user == nil {
		return "", "", nil, errcode.AuthError
	}
	// 2. 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, errcode.AuthError
	}
	// 3. 生成双Token
	accessToken, refreshToken, err = jwt.GenerateTokens(user.ID, user.Username, string(user.Role), string(user.Department))
	if err != nil {
		return "", "", nil, errcode.DBError
	}
	return accessToken, refreshToken, user, errcode.Success
}

// 刷新Token业务
func RefreshToken(refreshToken string) (newAccessToken, newRefreshToken string, errCode errcode.ErrCode) {
	// 1. 解析RefreshToken获取用户ID
	userID, errCode := jwt.ParseRefreshToken(refreshToken)
	if errCode != errcode.Success {
		return "", "", errCode
	}
	// 2. 查询用户
	user, err := dao.GetUserByID(userID)
	if err != nil || user == nil {
		return "", "", errcode.AuthError
	}
	// 3. 生成新的双Token
	newAccessToken, newRefreshToken, err = jwt.GenerateTokens(user.ID, user.Username, string(user.Role), string(user.Department))
	if err != nil {
		return "", "", errcode.DBError
	}
	return newAccessToken, newRefreshToken, errcode.Success
}

// 注销账号业务
func Logout(userID int64) errcode.ErrCode {
	if err := dao.DeleteUser(userID); err != nil {
		return errcode.DBError
	}
	return errcode.Success
}
func GetUserByID(userID int64) (*models.User, error) {
	return dao.GetUserByID(userID)
}
