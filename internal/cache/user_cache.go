package cache

import (
	"context"
	"fmt"
	"teamflow/internal/database"
	"teamflow/pkg/utils"
	"time"
)

// UserCache 用户缓存
type UserCache struct {
}

// GetUser 获取用户缓存中的用户，反序列化到 out。缓存不存在时返回 false, nil。
func (u *UserCache) GetUser(userID string, out interface{}) (bool, error) {
	return utils.GetCache(userID, out)
}

// SetUser 设置用户缓存，ttl<=0 表示不过期
func (u *UserCache) SetUser(userID string, user interface{}, ttl time.Duration) error {
	return utils.SetCache(userID, user, ttl)
}

// SetUserToken 设置用户token缓存
func (u *UserCache) SetUserToken(userID string, token string, ttl time.Duration) error {
	return utils.SetCache(fmt.Sprintf("%s:%s", userTokenPrefix, userID), token, ttl)
}
func (u *UserCache) GetUserToken(userID string) (string, error) {
	key := fmt.Sprintf("%s:%s", userTokenPrefix, userID)
	var token string
	ok, err := utils.GetCache(key, &token)
	if !ok || err != nil {
		return "", err
	}
	return token, nil
}
func (u *UserCache) DeleteUserToken(userID string) error {
	key := fmt.Sprintf("%s:%s", userTokenPrefix, userID)
	return database.RDB.Del(context.Background(), key).Err()
}
