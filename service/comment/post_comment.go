package comment

import (
	"fmt"
	"tiktok/logger"
	"tiktok/models"
)

const (
	CREATE = 1
	DELETE = 2
)

type PostCommentFlow struct {
	userId      int64
	videoId     int64
	commentId   int64
	actionType  int64
	commentText string
	comment     *models.Comment
	*Response
}

type Response struct {
	MyComment *models.Comment `json:"comment"`
}

// PostComment 发表评论 or 删除评论
func PostComment(userId int64, videoId int64, commentId int64, actionType int64, commentText string) (*Response, error) {
	return NewPostCommentFlow(userId, videoId, commentId, actionType, commentText).Do()
}

func NewPostCommentFlow(userId int64, videoId int64, commentId int64, actionType int64, commentText string) *PostCommentFlow {
	return &PostCommentFlow{
		userId:      userId,
		videoId:     videoId,
		commentId:   commentId,
		actionType:  actionType,
		commentText: commentText,
	}
}

func (p *PostCommentFlow) Do() (*Response, error) {
	if err := p.checkNum(); err != nil {
		logger.ZapLogger.Error("check num failed", logger.Error(err))
		return nil, err
	}
	if err := p.prepareData(); err != nil {
		logger.ZapLogger.Error("prepare data failed", logger.Error(err))
		return nil, err
	}
	return p.Response, nil
}

func (p *PostCommentFlow) prepareData() error {
	var err error
	switch p.actionType {
	case CREATE:
		p.comment, err = p.CreateComment()
	case DELETE:
		p.comment, err = p.DeleteComment()
	default:
		return fmt.Errorf("action type error")
	}
	return err
}

func (p *PostCommentFlow) CreateComment() (*models.Comment, error) {
	comment := &models.Comment{
		UserInfoId: p.userId,
		VideoId:    p.videoId,
		Content:    p.commentText,
	}
	err := models.NewCommentDAO().AddCommentAndUpdateCount(comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (p *PostCommentFlow) DeleteComment() (*models.Comment, error) {
	err := models.NewCommentDAO().DeleteCommentAndUpdateCountById(p.commentId, p.videoId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *PostCommentFlow) checkNum() error {
	if !models.NewUserInfoDAO().IsUserExistById(p.userId) {
		return fmt.Errorf("user not exist")
	}
	if !models.NewVideoDAO().IsVideoExistById(p.videoId) {
		return fmt.Errorf("video not exist")
	}
	if p.actionType == DELETE && !models.NewCommentDAO().IsCommentExistById(p.commentId) {
		return fmt.Errorf("comment not exist")
	}
	return nil
}
