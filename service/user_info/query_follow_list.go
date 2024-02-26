package user_info

import (
	"fmt"
	"tiktok/cache"
	"tiktok/logger"
	"tiktok/models"
)

type QueryFollowListFlow struct {
	userId   int64
	userList []*models.UserInfo
	*FollowList
}

type FollowList struct {
	UserList []*models.UserInfo `json:"user_list"`
}

// QueryFollowList 查询用户关注列表
func QueryFollowList(userId int64) (*FollowList, error) {
	return NewQueryFollowListFlow(userId).Do()
}

func NewQueryFollowListFlow(userId int64) *QueryFollowListFlow {
	return &QueryFollowListFlow{userId: userId}
}

func (q *QueryFollowListFlow) Do() (*FollowList, error) {
	if err := q.checkNum(); err != nil {
		logger.ZapLogger.Error("check num failed", logger.Error(err))
		return nil, err
	}

	if err := q.prepareData(); err != nil {
		logger.ZapLogger.Error("prepare data failed", logger.Error(err))
		return nil, err
	}

	q.FollowList = &FollowList{
		UserList: q.userList,
	}
	return q.FollowList, nil
}

func (q *QueryFollowListFlow) checkNum() error {
	if !models.NewUserInfoDAO().IsUserExistById(q.userId) {
		return fmt.Errorf("user not exist")
	}
	return nil
}

func (q *QueryFollowListFlow) prepareData() error {
	var userList []*models.UserInfo
	err := models.NewUserInfoDAO().GetFollowListByUserId(q.userId, &userList)
	if err != nil {
		return err
	}
	for _, v := range q.userList {
		v.IsFollow = cache.NewProxyIndexMap().GetUserRelation(q.userId, v.Id)
	}
	q.userList = userList
	return nil
}
