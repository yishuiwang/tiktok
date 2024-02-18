package user_info

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/models"
	"tiktok/service/user_info"
)

type FollowerListResponse struct {
	models.CommonResponse
	*user_info.FollowerList
}

func QueryFollowerHandler(c *gin.Context) {
	NewProxyQueryFollowerHandler(c).Do()
}

type ProxyQueryFollowerHandler struct {
	*gin.Context
	userId int64
	*user_info.FollowerList
}

func NewProxyQueryFollowerHandler(context *gin.Context) *ProxyQueryFollowerHandler {
	return &ProxyQueryFollowerHandler{Context: context}
}

func (p *ProxyQueryFollowerHandler) Do() {
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

func (p *ProxyQueryFollowerHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("parse user_id error")
	}
	p.userId = userId
	return nil
}

func (p *ProxyQueryFollowerHandler) prepareData() error {
	list, err := user_info.QueryFollowerList(p.userId)
	if err != nil {
		return err
	}
	p.FollowerList = list
	return nil
}

func (p *ProxyQueryFollowerHandler) SendError(msg string) {
	p.JSON(http.StatusOK, FollowerListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func (p *ProxyQueryFollowerHandler) SendOk(msg string) {
	p.JSON(http.StatusOK, FollowerListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		FollowerList: p.FollowerList,
	})
}
