package models

import (
	"gorm.io/gorm"
	"sync"
)

type UserInfo struct {
	Id            int64       `json:"id" gorm:"id,omitempty"`
	Name          string      `json:"name" gorm:"name,omitempty"`
	FollowCount   int64       `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount int64       `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow      bool        `json:"is_follow" gorm:"is_follow,omitempty"`
	User          *UserLogin  `json:"-"`                                     //用户与账号密码之间的一对一
	Videos        []*Video    `json:"-"`                                     //用户与投稿视频的一对多
	Follows       []*UserInfo `json:"-" gorm:"many2many:user_relations;"`    //用户之间的多对多
	FavorVideos   []*Video    `json:"-" gorm:"many2many:user_favor_videos;"` //用户与点赞视频之间的多对多
	Comments      []*Comment  `json:"-"`                                     //用户与评论的一对多
}

type UserInfoDAO struct {
}

var (
	userInfoDAO  *UserInfoDAO
	userInfoOnce sync.Once
)

func NewUserInfoDAO() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = &UserInfoDAO{}
	})
	return userInfoDAO
}

// QueryUserInfoById 根据用户id查询用户信息
func (u *UserInfoDAO) QueryUserInfoById(userId int64, userinfo *UserInfo) error {
	return DB.Where("id = ?", userId).First(userinfo).Error
}

// AddUserInfo 添加用户信息
func (u *UserInfoDAO) AddUserInfo(userinfo *UserInfo) error {
	return DB.Create(userinfo).Error
}

// IsUserExistById 根据用户id判断用户是否存在
func (u *UserInfoDAO) IsUserExistById(id int64) bool {
	var count int64
	DB.Model(&UserInfo{}).Where("id = ?", id).Count(&count)
	return count > 0
}

// AddUserFollow 添加用户关注
func (u *UserInfoDAO) AddUserFollow(userId, userToId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新关注数
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count+1 WHERE id = ?", userId).Error; err != nil {
			return err
		}
		// 2. 更新粉丝数
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count+1 WHERE id = ?", userToId).Error; err != nil {
			return err
		}
		// 3. 插入用户关注用户的记录
		if err := tx.Exec("INSERT INTO `user_relations` (`user_info_id`,`follow_id`) VALUES (?,?)", userId, userToId).Error; err != nil {
			return err
		}
		return nil
	})
}

// CancelUserFollow 取消用户关注
func (u *UserInfoDAO) CancelUserFollow(userId, userToId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新关注数
		if err := tx.Exec("UPDATE user_infos SET follow_count=follow_count-1 WHERE id = ? AND follow_count>0", userId).Error; err != nil {
			return err
		}
		// 2. 更新粉丝数
		if err := tx.Exec("UPDATE user_infos SET follower_count=follower_count-1 WHERE id = ? AND follower_count>0", userToId).Error; err != nil {
			return err
		}
		// 3. 删除用户关注用户的记录
		if err := tx.Exec("DELETE FROM `user_relations` WHERE user_info_id=? AND follow_id=?", userId, userToId).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetFollowListByUserId 根据用户id查询关注列表
func (u *UserInfoDAO) GetFollowListByUserId(userId int64, userList *[]*UserInfo) error {
	return DB.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.user_info_id = ? AND r.follow_id = u.id", userId).Scan(userList).Error
}

// GetFollowerListByUserId 根据用户id查询粉丝列表
func (u *UserInfoDAO) GetFollowerListByUserId(userId int64, userList *[]*UserInfo) error {
	return DB.Raw("SELECT u.* FROM user_relations r, user_infos u WHERE r.follow_id = ? AND r.user_info_id = u.id", userId).Scan(userList).Error
}
