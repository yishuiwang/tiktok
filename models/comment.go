package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	Id         int64     `json:"id"`
	UserInfoId int64     `json:"-"` //用于一对多关系的id
	VideoId    int64     `json:"-"` //一对多，视频对评论
	User       UserInfo  `json:"user" gorm:"-"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date" gorm:"-"`
}

type CommentDAO struct {
}

var (
	commentDao CommentDAO
)

func NewCommentDAO() *CommentDAO {
	return &commentDao
}

// AddCommentAndUpdateCount 添加评论并更新评论数
func (c *CommentDAO) AddCommentAndUpdateCount(comment *Comment) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 添加评论
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		// 2. 更新评论数
		if err := tx.Exec("UPDATE videos SET comment_count=comment_count+1 WHERE id = ?", comment.VideoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// DeleteCommentAndUpdateCountById 删除评论并更新评论数
func (c *CommentDAO) DeleteCommentAndUpdateCountById(commentId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除评论
		if err := tx.Exec("DELETE FROM comments WHERE id = ?", commentId).Error; err != nil {
			return err
		}
		// 2. 更新评论数
		if err := tx.Exec("UPDATE videos SET comment_count=comment_count-1 WHERE id = ?", videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// QueryCommentById 根据评论id查询评论
func (c *CommentDAO) QueryCommentById(id int64, comment *Comment) error {
	return DB.Where("id = ?", id).First(comment).Error
}

// QueryCommentListByVideoId 根据视频id查询评论列表
func (c *CommentDAO) QueryCommentListByVideoId(videoId int64, comments *[]*Comment) error {
	return DB.Model(&Comment{}).Where("video_id=?", videoId).Find(comments).Error
}

// IsCommentExistById 根据评论id判断评论是否存在
func (c *CommentDAO) IsCommentExistById(id int64) bool {
	var count int64
	DB.Model(&Comment{}).Where("id = ?", id).Count(&count)
	return count > 0
}
