package video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/models"
	"tiktok/service/video"
)

type ListResponse struct {
	models.CommonResponse
	*video.List
}

type ProxyQueryVideoList struct {
	c *gin.Context
}

func QueryVideoListHandler(c *gin.Context) {
	p := NewProxyQueryVideoList(c)
	rawId, _ := c.Get("user_id")
	err := p.DoQueryVideoListByUserId(rawId)
	if err != nil {
		p.QueryVideoListError(err.Error())
	}
}

func NewProxyQueryVideoList(c *gin.Context) *ProxyQueryVideoList {
	return &ProxyQueryVideoList{c: c}
}

func (p *ProxyQueryVideoList) DoQueryVideoListByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("userId parse error")
	}

	videoList, err := video.QueryVideoListByUserId(userId)
	if err != nil {
		return err
	}

	p.QueryVideoListOk(videoList)
	return nil
}

func (p *ProxyQueryVideoList) QueryVideoListError(msg string) {
	p.c.JSON(http.StatusOK, ListResponse{CommonResponse: models.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	}})
}

func (p *ProxyQueryVideoList) QueryVideoListOk(videoList *video.List) {
	p.c.JSON(http.StatusOK, ListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
		},
		List: videoList,
	})
}
