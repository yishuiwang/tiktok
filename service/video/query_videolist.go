package video

import (
	"fmt"
	"tiktok/models"
)

type QueryVideoListByUserIdFlow struct {
	userId    int64
	videos    []*models.Video
	videoList *List
}

type List struct {
	Videos []*models.Video `json:"video_list,omitempty"`
}

// QueryVideoListByUserId 查询用户视频列表
func QueryVideoListByUserId(userId int64) (*List, error) {
	return NewQueryVideoListByUserIdFlow(userId).Do()
}

func NewQueryVideoListByUserIdFlow(userId int64) *QueryVideoListByUserIdFlow {
	return &QueryVideoListByUserIdFlow{userId: userId}
}

func (q *QueryVideoListByUserIdFlow) Do() (*List, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

func (q *QueryVideoListByUserIdFlow) checkNum() error {
	if !models.NewUserInfoDAO().IsUserExistById(q.userId) {
		return fmt.Errorf("user not exist")
	}
	return nil
}

func (q *QueryVideoListByUserIdFlow) packData() error {
	var videos []*models.Video
	err := models.NewVideoDAO().QueryVideoListByUserId(q.userId, &videos)
	if err != nil {
		return err
	}
	q.videoList = &List{
		Videos: videos,
	}
	return nil
}
