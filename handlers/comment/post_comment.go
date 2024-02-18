package comment

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/models"
	"tiktok/service/comment"
)

type PostCommentResponse struct {
	models.CommonResponse
	*comment.Response
}

type ProxyPostCommentHandler struct {
	*gin.Context
	videoId     int64
	userId      int64
	commentId   int64
	actionType  int64
	commentText string
}

func PostCommentHandler(c *gin.Context) {
	NewProxyPostCommentHandler(c).Do()
}

func NewProxyPostCommentHandler(context *gin.Context) *ProxyPostCommentHandler {
	return &ProxyPostCommentHandler{Context: context}
}

func (p *ProxyPostCommentHandler) Do() {
	if err := p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	commentRes, err := comment.PostComment(p.userId, p.videoId, p.commentId, p.actionType, p.commentText)
	if err != nil {
		p.SendError(err.Error())
		return
	}
	p.SendOk(commentRes)
}

func (p *ProxyPostCommentHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("parse user_id error")
	}
	p.userId = userId

	rawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		return err
	}
	p.videoId = videoId

	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	switch actionType {
	case comment.CREATE:
		p.commentText = p.Query("comment_text")
	case comment.DELETE:
		p.commentId, err = strconv.ParseInt(p.Query("comment_id"), 10, 64)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid action_type: %d", actionType)
	}
	p.actionType = actionType
	return nil
}

func (p *ProxyPostCommentHandler) SendError(msg string) {
	p.JSON(http.StatusOK, PostCommentResponse{
		CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: msg}, Response: &comment.Response{}})
}

func (p *ProxyPostCommentHandler) SendOk(comment *comment.Response) {
	p.JSON(http.StatusOK, PostCommentResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		Response:       comment,
	})
}
