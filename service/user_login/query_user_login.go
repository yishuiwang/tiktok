package user_login

import (
	"fmt"
	"tiktok/logger"
	"tiktok/middleware"
	"tiktok/models"
)

type QueryUserLoginFlow struct {
	username string
	password string

	data   *LoginResponse
	userid int64
	token  string
}

// QueryUserLogin 查询用户是否存在，并返回token和id
func QueryUserLogin(username, password string) (*LoginResponse, error) {
	return NewQueryUserLoginFlow(username, password).Do()
}

func NewQueryUserLoginFlow(username, password string) *QueryUserLoginFlow {
	return &QueryUserLoginFlow{username: username, password: password}
}

func (q *QueryUserLoginFlow) Do() (*LoginResponse, error) {
	if err := q.checkNum(); err != nil {
		logger.ZapLogger.Error("checkNum failed", logger.Error(err))
		return nil, err
	}

	if err := q.prepareData(); err != nil {
		logger.ZapLogger.Error("prepareData failed", logger.Error(err))
		return nil, err
	}

	q.data = &LoginResponse{
		UserId: q.userid,
		Token:  q.token,
	}

	return q.data, nil
}

func (q *QueryUserLoginFlow) checkNum() error {
	if q.username == "" || q.password == "" {
		return fmt.Errorf("username or password is empty")
	}
	if len(q.username) < 6 || len(q.username) > 20 {
		return fmt.Errorf("username length is not valid")
	}
	if len(q.password) < 6 || len(q.password) > 20 {
		return fmt.Errorf("password length is not valid")
	}
	return nil
}

func (q *QueryUserLoginFlow) prepareData() error {
	var login models.UserLogin
	if err := models.NewUserLoginDAO().QueryUserLogin(q.username, q.password, &login); err != nil {
		return err
	}
	q.userid = login.UserInfoId

	token, err := middleware.ReleaseToken(login)
	if err != nil {
		return err
	}
	q.token = token
	return nil
}
