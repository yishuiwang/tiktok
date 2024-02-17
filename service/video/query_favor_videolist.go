package video

import (
	"fmt"
	"tiktok/models"
)

type QueryFavorVideoListFlow struct {
	userId    int64
	videos    []*models.Video
	videoList *FavorList
}

type FavorList struct {
	Videos []*models.Video `json:"video_list"`
}

// QueryFavorVideoList 查询用户收藏视频列表
func QueryFavorVideoList(userId int64) (*FavorList, error) {
	return NewQueryFavorVideoListFlow(userId).Do()
}

func NewQueryFavorVideoListFlow(userId int64) *QueryFavorVideoListFlow {
	return &QueryFavorVideoListFlow{userId: userId}
}

func (q *QueryFavorVideoListFlow) Do() (*FavorList, error) {
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

func (q *QueryFavorVideoListFlow) checkNum() error {
	if !models.NewUserInfoDAO().IsUserExistById(q.userId) {
		return fmt.Errorf("user not exist")
	}
	return nil
}

func (q *QueryFavorVideoListFlow) prepareData() error {
	var videos []*models.Video
	err := models.NewVideoDAO().QueryFavorVideoListByUserId(q.userId, &videos)
	if err != nil {
		return err
	}
	q.videoList = &FavorList{
		Videos: videos,
	}
	return nil
}
