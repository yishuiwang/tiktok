package models

import "sync"

// UserLogin 用户登录表，和UserInfo属于一对一关系
type UserLogin struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

type UserLoginDAO struct {
}

var (
	userLoginDao  *UserLoginDAO
	userLoginOnce sync.Once
)

func NewUserLoginDAO() *UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDao = &UserLoginDAO{}
	})
	return userLoginDao
}

// QueryUserLogin 根据用户名和密码查询用户登录信息
func (u *UserLoginDAO) QueryUserLogin(username, password string, login *UserLogin) error {
	return DB.Where("username = ? AND password = ?", username, password).First(login).Error
}

// IsUserExistByUsername 根据用户名判断用户是否存在
func (u *UserLoginDAO) IsUserExistByUsername(username string) bool {
	var count int64
	DB.Model(&UserLogin{}).Where("username = ?", username).Count(&count)
	return count > 0
}
