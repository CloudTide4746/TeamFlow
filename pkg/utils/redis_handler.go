package utils

import (
	"context"
	"encoding/json"
	"time"

	"teamflow/internal/database"

	"github.com/redis/go-redis/v9"
)

// redisOpTimeout 单次 Redis 操作的超时时间
const redisOpTimeout = 2 * time.Second

// SetCache 将 value 序列化为 JSON 后写入 Redis。ttl<=0 表示不过期。
func SetCache(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), redisOpTimeout)
	defer cancel()
	return database.RDB.Set(ctx, key, string(data), ttl).Err()
}

// GetCache 读取 Redis 中的 JSON 值并反序列化到 out。key 不存在时返回 false, nil。
func GetCache(key string, out interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), redisOpTimeout)
	defer cancel()
	val, err := database.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(val), out); err != nil {
		return false, err
	}
	return true, nil
}
