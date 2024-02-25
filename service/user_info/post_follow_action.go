package user_info

import (
	"fmt"
	"tiktok/logger"
	"tiktok/models"
)

const (
	FOLLOW = 1
	CANCEL = 2
)

type PostFollowActionFlow struct {
	userId     int64
	userToId   int64
	actionType int
}

// PostFollowAction 关注用户
func PostFollowAction(userId, userToId int64, actionType int) error {
	return NewPostFollowActionFlow(userId, userToId, actionType).Do()
}

func NewPostFollowActionFlow(userId, userToId int64, actionType int) *PostFollowActionFlow {
	return &PostFollowActionFlow{userId: userId, userToId: userToId, actionType: actionType}
}

func (p *PostFollowActionFlow) Do() error {
	if err := p.checkNum(); err != nil {
		logger.ZapLogger.Error("check num failed", logger.Error(err))
		return err
	}

	if err := p.publish(); err != nil {
		logger.ZapLogger.Error("publish failed", logger.Error(err))
		return err
	}

	return nil
}

func (p *PostFollowActionFlow) checkNum() error {
	if p.actionType != FOLLOW && p.actionType != CANCEL {
		return fmt.Errorf("action type is not valid")
	}

	if p.userId == p.userToId {
		return fmt.Errorf("can not follow yourself")
	}

	if models.NewUserInfoDAO().IsUserExistById(p.userToId) == false {
		return fmt.Errorf("user to follow not exist")
	}

	if models.NewUserInfoDAO().IsUserExistById(p.userId) == false {
		return fmt.Errorf("user not exist")
	}

	return nil
}

func (p *PostFollowActionFlow) publish() error {
	userDAO := models.NewUserInfoDAO()
	var err error
	switch p.actionType {
	case FOLLOW:
		err = userDAO.AddUserFollow(p.userId, p.userToId)
		//更新redis的关注信息
		//cache.NewProxyIndexMap().UpdateUserRelation(p.userId, p.userToId, true)
	case CANCEL:
		err = userDAO.CancelUserFollow(p.userId, p.userToId)
		//cache.NewProxyIndexMap().UpdateUserRelation(p.userId, p.userToId, false)
	default:
		return fmt.Errorf("action type is not valid")
	}
	return err
}
