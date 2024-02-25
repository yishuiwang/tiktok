package video

import (
	"fmt"
	"tiktok/logger"
	"tiktok/models"
)

const (
	PLUS  = 1
	MINUS = 2
)

type PostFavorStateFlow struct {
	userId     int64
	videoId    int64
	actionType int64
}

func PostFavorState(userId, videoId, actionType int64) error {
	return NewPostFavorStateFlow(userId, videoId, actionType).Do()
}

func NewPostFavorStateFlow(userId, videoId, actionType int64) *PostFavorStateFlow {
	return &PostFavorStateFlow{userId: userId, videoId: videoId, actionType: actionType}
}

func (p *PostFavorStateFlow) Do() error {
	if err := p.checkNum(); err != nil {
		logger.ZapLogger.Error("PostFavorStateFlow checkNum failed", logger.Error(err))
		return err
	}

	switch p.actionType {
	case PLUS:
		return p.PlusOperation()
	case MINUS:
		return p.MinusOperation()
	default:
		return fmt.Errorf("action type is not valid")
	}
}

// PlusOperation 点赞操作
func (p *PostFavorStateFlow) PlusOperation() error {
	//视频点赞数目+1
	err := models.NewVideoDAO().PlusOneFavorByUserIdAndVideoId(p.userId, p.videoId)
	if err != nil {
		return fmt.Errorf("plus one favor failed")
	}
	//对应的用户是否点赞的映射状态更新
	//cache.NewProxyIndexMap().UpdateVideoFavorState(p.userId, p.videoId, true)
	return nil
}

// MinusOperation 取消点赞操作
func (p *PostFavorStateFlow) MinusOperation() error {
	//视频点赞数目-1
	err := models.NewVideoDAO().MinusOneFavorByUserIdAndVideoId(p.userId, p.videoId)
	if err != nil {
		return fmt.Errorf("minus one favor failed")
	}
	//对应的用户是否点赞的映射状态更新
	//cache.NewProxyIndexMap().UpdateVideoFavorState(p.userId, p.videoId, false)
	return nil
}

func (p *PostFavorStateFlow) checkNum() error {
	if p.actionType != PLUS && p.actionType != MINUS {
		return fmt.Errorf("action type is not valid")
	}
	if models.NewUserInfoDAO().IsUserExistById(p.userId) == false {
		return fmt.Errorf("user not exist")
	}
	return nil
}
