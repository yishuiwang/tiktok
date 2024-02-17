package user_info

import (
	"fmt"
	"tiktok/models"
)

type FollowerList struct {
	UserList []*models.UserInfo `json:"user_list"`
}

type QueryFollowerListFlow struct {
	userId   int64
	userList []*models.UserInfo
	*FollowerList
}

// QueryFollowerList 查询用户粉丝列表
func QueryFollowerList(userId int64) (*FollowerList, error) {
	return NewQueryFollowerListFlow(userId).Do()
}

func NewQueryFollowerListFlow(userId int64) *QueryFollowerListFlow {
	return &QueryFollowerListFlow{userId: userId}
}

func (q *QueryFollowerListFlow) Do() (*FollowerList, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}

	if err := q.prepareData(); err != nil {
		return nil, err
	}

	q.FollowerList = &FollowerList{
		UserList: q.userList,
	}
	return q.FollowerList, nil
}

func (q *QueryFollowerListFlow) checkNum() error {
	if !models.NewUserInfoDAO().IsUserExistById(q.userId) {
		return fmt.Errorf("user not exist")
	}
	return nil
}

func (q *QueryFollowerListFlow) prepareData() error {
	var userList []*models.UserInfo
	err := models.NewUserInfoDAO().GetFollowerListByUserId(q.userId, &userList)
	if err != nil {
		return err
	}

	//for _, v := range q.userList {
	//v.IsFollow = cache.NewProxyIndexMap().GetUserRelation(q.userId, v.Id)
	//}
	return nil
}
