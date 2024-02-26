package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"tiktok/config"
)

var ctx = context.Background()
var rdb *redis.Client

const (
	favor    = "favor"
	relation = "relation"
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.GetConfig("redis.addr"),
		Password: config.GetConfig("redis.password"),
		DB:       0, // use default DB
	})
}

var (
	proxyIndexOperation ProxyIndexMap
)

type ProxyIndexMap struct {
}

func NewProxyIndexMap() *ProxyIndexMap {
	return &proxyIndexOperation
}

// UpdateVideoFavorState 更新视频点赞状态，true为点赞，false为取消点赞
func (i *ProxyIndexMap) UpdateVideoFavorState(userId int64, videoId int64, state bool) {
	key := fmt.Sprintf("%s:%d", favor, userId)
	if state {
		// 如果点赞，就加入集合
		rdb.SAdd(ctx, key, videoId)
	} else {
		// 如果取消点赞，就从集合中删除
		rdb.SRem(ctx, key, videoId)
	}
}

// GetVideoFavorState 获取视频点赞状态
func (i *ProxyIndexMap) GetVideoFavorState(userId int64, videoId int64) bool {
	key := fmt.Sprintf("%s:%d", favor, userId)
	// 判断是否存在
	return rdb.SIsMember(ctx, key, videoId).Val()
}

// UpdateUserRelation 更新用户关系，true为关注，false为取消关注
func (i *ProxyIndexMap) UpdateUserRelation(userId int64, targetUserId int64, state bool) {
	key := fmt.Sprintf("%s:%d", relation, userId)
	if state {
		// 如果关注，就加入集合
		rdb.SAdd(ctx, key, targetUserId)
	} else {
		// 如果取消关注，就从集合中删除
		rdb.SRem(ctx, key, targetUserId)
	}
}

// GetUserRelation 获取用户关系
func (i *ProxyIndexMap) GetUserRelation(userId int64, targetUserId int64) bool {
	key := fmt.Sprintf("%s:%d", relation, userId)
	// 判断是否存在
	return rdb.SIsMember(ctx, key, targetUserId).Val()
}
