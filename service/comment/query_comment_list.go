package comment

import (
	"fmt"
	"tiktok/logger"
	"tiktok/models"
)

type QueryCommentListFlow struct {
	userId      int64
	videoId     int64
	comments    []*models.Comment
	commentList *List
}

type List struct {
	Comments []*models.Comment `json:"comment_list"`
}

// QueryCommentList 查询评论列表
func QueryCommentList(userId, videoId int64) (*List, error) {
	return NewQueryCommentListFlow(userId, videoId).Do()
}

func NewQueryCommentListFlow(userId, videoId int64) *QueryCommentListFlow {
	return &QueryCommentListFlow{
		userId:  userId,
		videoId: videoId,
	}
}

func (q *QueryCommentListFlow) Do() (*List, error) {
	if err := q.checkNum(); err != nil {
		logger.ZapLogger.Error("check num failed", logger.Error(err))
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		logger.ZapLogger.Error("prepare data failed", logger.Error(err))
		return nil, err
	}
	return q.commentList, nil
}

func (q *QueryCommentListFlow) checkNum() error {
	if !models.NewUserInfoDAO().IsUserExistById(q.userId) {
		return fmt.Errorf("user not exist")
	}
	if !models.NewVideoDAO().IsVideoExistById(q.videoId) {
		return fmt.Errorf("video not exist")
	}
	return nil
}

func (q *QueryCommentListFlow) prepareData() error {
	var comments []*models.Comment
	err := models.NewCommentDAO().QueryCommentListByVideoId(q.videoId, &comments)
	if err != nil {
		return err
	}
	q.commentList = &List{
		Comments: comments,
	}
	return nil
}
