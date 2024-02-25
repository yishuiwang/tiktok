package user_login

import (
	"fmt"
	"tiktok/logger"
	"tiktok/middleware"
	"tiktok/models"
)

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type PostUserLoginFlow struct {
	username string
	password string

	data   *LoginResponse
	userid int64
	token  string
}

// PostUserLogin 注册用户并得到token和id
func PostUserLogin(username, password string) (*LoginResponse, error) {
	return NewPostUserLoginFlow(username, password).Do()
}

func NewPostUserLoginFlow(username, password string) *PostUserLoginFlow {
	return &PostUserLoginFlow{username: username, password: password}
}

func (q *PostUserLoginFlow) Do() (*LoginResponse, error) {
	if err := q.checkNum(); err != nil {
		logger.ZapLogger.Error("checkNum failed", logger.Error(err))
		return nil, err
	}

	if err := q.updateData(); err != nil {
		logger.ZapLogger.Error("updateData failed", logger.Error(err))
		return nil, err
	}

	q.data = &LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}

	return q.data, nil
}

func (q *PostUserLoginFlow) checkNum() error {
	if len(q.username) == 0 || len(q.password) == 0 {
		return fmt.Errorf("username or password is empty")
	}

	if len(q.username) < 6 || len(q.username) > 20 {
		return fmt.Errorf("username length should be 6-20")
	}

	if len(q.password) < 6 || len(q.password) > 20 {
		return fmt.Errorf("password length should be 6-20")
	}

	return nil
}

func (q *PostUserLoginFlow) updateData() error {
	userLogin := models.UserLogin{
		Username: q.username,
		Password: q.password,
	}
	userInfo := models.UserInfo{
		User: &userLogin,
		Name: q.username,
	}

	if models.NewUserLoginDAO().IsUserExistByUsername(q.username) {
		return fmt.Errorf("username already exists")
	}

	if err := models.NewUserInfoDAO().AddUserInfo(&userInfo); err != nil {
		return err
	}

	token, err := middleware.ReleaseToken(userLogin)
	if err != nil {
		return err
	}

	q.token = token
	q.userid = userLogin.Id

	return nil
}
