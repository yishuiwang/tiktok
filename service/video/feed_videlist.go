package video

import (
	"tiktok/logger"
	"tiktok/models"
	"time"
)

const MaxVideoNum = 10

type FeedVideoList struct {
	Videos   []*models.Video `json:"video_list,omitempty"`
	NextTime int64           `json:"next_time,omitempty"`
}

type QueryFeedVideoListFlow struct {
	userId     int64
	latestTime time.Time
	videos     []*models.Video
	nextTime   int64
	feedVideo  *FeedVideoList
}

// QueryFeedVideoList 查询用户关注的视频列表
func QueryFeedVideoList(userId int64, latestTime time.Time) (*FeedVideoList, error) {
	return NewQueryFeedVideoListFlow(userId, latestTime).Do()
}

func NewQueryFeedVideoListFlow(userId int64, latestTime time.Time) *QueryFeedVideoListFlow {
	return &QueryFeedVideoListFlow{userId: userId, latestTime: latestTime}
}

func (q *QueryFeedVideoListFlow) Do() (*FeedVideoList, error) {
	if err := q.prepareData(); err != nil {
		logger.ZapLogger.Error("QueryFeedVideoListFlow prepareData failed", logger.Error(err))
		return nil, err
	}

	q.feedVideo = &FeedVideoList{
		Videos:   q.videos,
		NextTime: q.nextTime,
	}
	return q.feedVideo, nil
}

func (q *QueryFeedVideoListFlow) prepareData() error {
	var videos []*models.Video
	err := models.NewVideoDAO().QueryVideoListByLimitAndTime(MaxVideoNum, q.latestTime, &q.videos)
	if err != nil {
		return err
	}
	q.videos = videos
	return nil
}
