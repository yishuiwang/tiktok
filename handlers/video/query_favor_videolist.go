package video

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/models"
	"tiktok/service/video"
)

type FavorVideoListResponse struct {
	models.CommonResponse
	*video.FavorList
}

type ProxyFavorVideoListHandler struct {
	*gin.Context
	userId int64
}

func QueryFavorVideoListHandler(c *gin.Context) {
	NewProxyFavorVideoListHandler(c).Do()
}

func NewProxyFavorVideoListHandler(c *gin.Context) *ProxyFavorVideoListHandler {
	return &ProxyFavorVideoListHandler{Context: c}
}

func (p *ProxyFavorVideoListHandler) Do() {
	if err := p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}

	favorVideoList, err := video.QueryFavorVideoList(p.userId)
	if err != nil {
		p.SendError(err.Error())
		return
	}

	p.SendOk(favorVideoList)
}

func (p *ProxyFavorVideoListHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId parse error")
	}
	p.userId = userId
	return nil
}

func (p *ProxyFavorVideoListHandler) SendError(msg string) {
	p.JSON(http.StatusOK, FavorVideoListResponse{
		CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: msg}})
}

func (p *ProxyFavorVideoListHandler) SendOk(favorList *video.FavorList) {
	p.JSON(http.StatusOK, FavorVideoListResponse{CommonResponse: models.CommonResponse{StatusCode: 0},
		FavorList: favorList,
	})
}
