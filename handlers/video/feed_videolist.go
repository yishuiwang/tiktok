package video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/middleware"
	"tiktok/models"
	"tiktok/service/video"
	"time"
)

type FeedResponse struct {
	models.CommonResponse
	*video.FeedVideoList
}

type ProxyFeedVideoList struct {
	*gin.Context
}

func FeedVideoListHandler(c *gin.Context) {
	p := NewProxyFeedVideoList(c)
	token, ok := c.GetQuery("token")
	if !ok {
		err := p.DoNoToken()
		if err != nil {
			p.FeedVideoListError(err.Error())
		}
		return
	}

	err := p.DoHasToken(token)
	if err != nil {
		p.FeedVideoListError(err.Error())
	}
}

func NewProxyFeedVideoList(c *gin.Context) *ProxyFeedVideoList {
	return &ProxyFeedVideoList{Context: c}
}

func (p *ProxyFeedVideoList) DoNoToken() error {
	rawTimestamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, intTime*1e6)
	}
	videoList, err := video.QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListOk(videoList)
	return nil
}

func (p *ProxyFeedVideoList) DoHasToken(token string) error {
	if claim, ok := middleware.ParseToken(token); ok {
		if time.Now().Unix() > claim.ExpiresAt {
			return errors.New("token timeout")
		}
		rawTimestamp := p.Query("latest_time")
		var latestTime time.Time
		intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
		if err != nil {
			latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
		}
		videoList, err := video.QueryFeedVideoList(claim.UserId, latestTime)
		if err != nil {
			return err
		}
		p.FeedVideoListOk(videoList)
		return nil
	}
	return errors.New("token error")
}

func (p *ProxyFeedVideoList) FeedVideoListError(msg string) {
	p.JSON(http.StatusOK, FeedResponse{CommonResponse: models.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	}})
}

func (p *ProxyFeedVideoList) FeedVideoListOk(videoList *video.FeedVideoList) {
	p.JSON(http.StatusOK, FeedResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
		},
		FeedVideoList: videoList,
	},
	)
}
