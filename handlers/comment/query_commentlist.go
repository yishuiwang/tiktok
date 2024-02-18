package comment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/models"
	"tiktok/service/comment"
)

type ListResponse struct {
	models.CommonResponse
	*comment.List
}

type ProxyCommentListHandler struct {
	*gin.Context
	videoId int64
	userId  int64
}

func QueryCommentListHandler(c *gin.Context) {
	NewProxyCommentListHandler(c).Do()
}

func NewProxyCommentListHandler(context *gin.Context) *ProxyCommentListHandler {
	return &ProxyCommentListHandler{Context: context}
}

func (p *ProxyCommentListHandler) Do() {
	if err := p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}

	commentList, err := comment.QueryCommentList(p.userId, p.videoId)
	if err != nil {
		p.SendError(err.Error())
		return
	}

	p.SendOk(commentList)
}

func (p *ProxyCommentListHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId

	rawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		return err
	}
	p.videoId = videoId

	return nil
}

func (p *ProxyCommentListHandler) SendError(msg string) {
	p.JSON(http.StatusOK, ListResponse{CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: msg}})
}

func (p *ProxyCommentListHandler) SendOk(commentList *comment.List) {
	p.JSON(http.StatusOK, ListResponse{CommonResponse: models.CommonResponse{StatusCode: 0},
		List: commentList,
	})
}
