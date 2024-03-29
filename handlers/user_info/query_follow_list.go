package user_info

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/models"
	"tiktok/service/user_info"
)

type FollowListResponse struct {
	models.CommonResponse
	*user_info.FollowList
}

type ProxyQueryFollowList struct {
	*gin.Context
	userId int64
	*user_info.FollowList
}

func QueryFollowListHandler(c *gin.Context) {
	NewProxyQueryFollowList(c).Do()
}

func NewProxyQueryFollowList(context *gin.Context) *ProxyQueryFollowList {
	return &ProxyQueryFollowList{Context: context}
}

func (p *ProxyQueryFollowList) Do() {
	var err error
	if err = p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	if err = p.prepareData(); err != nil {
		p.SendError(err.Error())
		return
	}
	p.SendOk("success")
}

func (p *ProxyQueryFollowList) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("user_id is not found")
	}
	p.userId = userId
	return nil
}

func (p *ProxyQueryFollowList) prepareData() error {
	list, err := user_info.QueryFollowList(p.userId)
	if err != nil {
		return err
	}
	p.FollowList = list
	return nil
}

func (p *ProxyQueryFollowList) SendError(msg string) {
	p.JSON(http.StatusOK, FollowListResponse{
		CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *ProxyQueryFollowList) SendOk(msg string) {
	p.JSON(http.StatusOK, FollowListResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0, StatusMsg: msg},
		FollowList:     p.FollowList,
	})
}
