package user_info

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/models"
)

type UserResponse struct {
	models.CommonResponse
	User *models.UserInfo `json:"user,omitempty"`
}

type ProxyUserInfo struct {
	c *gin.Context
}

func UserInfoHandler(c *gin.Context) {
	p := NewProxyUserInfo(c)
	rawId, ok := c.Get("user_id")
	if !ok {
		p.UserInfoError("user_id not found")
		return
	}
	err := p.DoQueryUserInfoByUserId(rawId)
	if err != nil {
		p.UserInfoError(err.Error())
	}
}

func NewProxyUserInfo(c *gin.Context) *ProxyUserInfo {
	return &ProxyUserInfo{c: c}
}

func (p *ProxyUserInfo) DoQueryUserInfoByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("parse userId error")
	}
	userinfoDAO := models.NewUserInfoDAO()

	var userInfo models.UserInfo
	err := userinfoDAO.QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOk(&userInfo)
	return nil
}

func (p *ProxyUserInfo) UserInfoError(msg string) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *ProxyUserInfo) UserInfoOk(user *models.UserInfo) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		User:           user,
	})
}
