package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"west.garden/template/common/config"
	"west.garden/template/common/consts"
	"strconv"
)

// RedisClient
var RedisClient *redis.Client

// Init
func Init(redisConfig *config.RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password:     redisConfig.Auth,
		DB:           redisConfig.Db,
		PoolSize:     redisConfig.MaxActive,
		MinIdleConns: redisConfig.MaxIdle,
	})
	return nil
}

var ctx = context.Background()

func AddScore(setName consts.RedisStoreKey, id int64, initialScore float64) error {
	idStr := strconv.FormatInt(id, 10)
	return RedisClient.ZAdd(ctx, string(setName), &redis.Z{Score: initialScore, Member: idStr}).Err()
}

func GetScoreRankById(setName consts.RedisStoreKey, id int64) (int64, error) {
	idStr := strconv.FormatInt(id, 10)
	rank, err := RedisClient.ZRevRank(ctx, string(setName), idStr).Result()
	if err != nil {
		return -1, err
	}
	return rank + 1, nil //rank from 0, so +1
}

func DeleteScore(setName consts.RedisStoreKey, id int64) error {
	idStr := strconv.FormatInt(id, 10)
	return RedisClient.ZRem(ctx, string(setName), idStr).Err()
}

func GetScoreById(setName consts.RedisStoreKey, id int64) (float64, error) {
	idStr := strconv.FormatInt(id, 10)
	score, err := RedisClient.ZScore(ctx, string(setName), idStr).Result()
	if err != nil {
		return 0, err
	}
	return score, nil
}

func DeleteKey(setName consts.RedisStoreKey) error {
	return RedisClient.Del(ctx, string(setName)).Err()
}
