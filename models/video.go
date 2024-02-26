package models

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Video struct {
	Id            int64       `json:"id,omitempty"`
	UserInfoId    int64       `json:"-"`
	Author        UserInfo    `json:"author,omitempty" gorm:"-"` //这里应该是作者对视频的一对多的关系，而不是视频对作者，故gorm不能存他，但json需要返回它
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"`
	Title         string      `json:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	Comments      []*Comment  `json:"-"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}

type VideoDAO struct {
}

var (
	videoDAO  *VideoDAO
	videoOnce sync.Once
)

func NewVideoDAO() *VideoDAO {
	videoOnce.Do(func() {
		videoDAO = &VideoDAO{}
	})
	return videoDAO
}

// AddVideo 添加视频
func (v *VideoDAO) AddVideo(video *Video) error {
	return DB.Create(video).Error
}

// QueryVideoByVideoId 根据视频id查询视频
func (v *VideoDAO) QueryVideoByVideoId(videoId int64, video *Video) error {
	return DB.Where("id = ?", videoId).First(video).Error
}

// QueryVideoCountByUserId 根据用户id查询视频
func (v *VideoDAO) QueryVideoCountByUserId(userId int64, count *int64) error {
	return DB.Model(&Video{}).Where("user_info_id = ?", userId).Count(count).Error
}

// QueryVideoListByUserId 根据用户id查询视频列表
func (v *VideoDAO) QueryVideoListByUserId(userId int64, videoList *[]*Video) error {
	return DB.Where("user_info_id = ?", userId).Find(videoList).Error
}

// QueryVideoListByLimitAndTime 根据数量和时间查询视频列表
func (v *VideoDAO) QueryVideoListByLimitAndTime(limit int, latestTime time.Time, videoList *[]*Video) error {
	return DB.Where("created_at < ?", latestTime).Order("created_at desc").Limit(limit).Find(videoList).Error
}

// PlusOneFavorByUserIdAndVideoId 用户对视频点赞
func (v *VideoDAO) PlusOneFavorByUserIdAndVideoId(userId int64, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新视频的点赞数
		if err := tx.Exec("UPDATE videos SET favorite_count=favorite_count+1 WHERE id = ?", videoId).Error; err != nil {
			return err
		}
		// 2. 插入用户点赞视频的记录
		if err := tx.Exec("INSERT INTO `user_favor_videos` (`user_info_id`,`video_id`) VALUES (?,?)", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// MinusOneFavorByUserIdAndVideoId 用户取消对视频点赞
func (v *VideoDAO) MinusOneFavorByUserIdAndVideoId(userId int64, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新视频的点赞数, 点赞数不能为负数
		if err := tx.Exec("UPDATE videos SET favorite_count=favorite_count-1 WHERE id = ? AND favorite_count>0", videoId).Error; err != nil {
			return err
		}
		// 2. 删除用户点赞视频的记录
		if err := tx.Exec("DELETE FROM `user_favor_videos` WHERE `user_info_id` = ? AND `video_id` = ?", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// QueryFavorVideoListByUserId 查询用户点赞的视频列表
func (v *VideoDAO) QueryFavorVideoListByUserId(userId int64, videoList *[]*Video) error {
	//多表查询，左连接得到结果，再映射到数据
	if err := DB.Model(&UserInfo{}).Where("id = ?", userId).Preload("Users").Find(&UserInfo{}).Error; err != nil {
		return err
	}

	//如果id为0，则说明没有查到数据
	if len(*videoList) == 0 || (*videoList)[0].Id == 0 {
		return fmt.Errorf("点赞列表为空")
	}
	return nil
}

func (v *VideoDAO) IsVideoExistById(videoId int64) bool {
	var count int64
	DB.Model(&Video{}).Where("id = ?", videoId).Count(&count)
	return count > 0
}
