package user_info

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/models"
	"tiktok/service/user_info"
)

func PostFollowActionHandler(c *gin.Context) {
	NewProxyPostFollowAction(c).Do()
}

type ProxyPostFollowAction struct {
	*gin.Context
	userId     int64
	followId   int64
	actionType int
}

func NewProxyPostFollowAction(context *gin.Context) *ProxyPostFollowAction {
	return &ProxyPostFollowAction{Context: context}
}

func (p *ProxyPostFollowAction) Do() {
	var err error
	if err = p.prepareNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	if err = p.startAction(); err != nil {
		p.SendError(err.Error())
		return
	}
	p.SendOk("操作成功")
}

func (p *ProxyPostFollowAction) prepareNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("user_id is not found")
	}
	p.userId = userId

	//解析需要关注的id
	followId := p.Query("to_user_id")
	parseInt, err := strconv.ParseInt(followId, 10, 64)
	if err != nil {
		return err
	}
	p.followId = parseInt

	//解析action_type
	actionType := p.Query("action_type")
	parseInt, err = strconv.ParseInt(actionType, 10, 32)
	if err != nil {
		return err
	}
	p.actionType = int(parseInt)
	return nil
}

func (p *ProxyPostFollowAction) startAction() error {
	err := user_info.PostFollowAction(p.userId, p.followId, p.actionType)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProxyPostFollowAction) SendError(msg string) {
	p.JSON(http.StatusOK, models.CommonResponse{StatusCode: 1, StatusMsg: msg})
}

func (p *ProxyPostFollowAction) SendOk(msg string) {
	p.JSON(http.StatusOK, models.CommonResponse{StatusCode: 1, StatusMsg: msg})
}
