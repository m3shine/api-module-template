package consts

import (
	"errors"
	"time"
)

type LockKey string

const (
	//RedisLockTtl 过期时间
	RedisLockTtl = time.Second * 30
	//RedisLockResetTTLInterval 重置过期时间间隔
	RedisLockResetTTLInterval = RedisLockTtl / 3
	//RedisLockTryLockInterval 重新获取锁间隔
	RedisLockTryLockInterval = time.Second
	//RedisLockUnlockScript 解锁脚本
	RedisLockUnlockScript = `
			if redis.call("get",KEYS[1]) == ARGV[1] then
				return redis.call("del",KEYS[1])
			else
				return 0
			end`

	//RedisKeyCPUAssemble CPU组装key
	RedisKeyCPUAssemble LockKey = "sqrd:cpu:assemble:%s"
	RedisKeyGPUAssemble LockKey = "sqrd:gpu:assemble:%s"
	RedisKeyOpenBox     LockKey = "sqrd:open:box:%s"
	RedisKeyClaimInitRewards LockKey = "sqrd:claim:init:rewards:%s"
)

var (
	// RedisLockLockFailedErr 加锁失败
	RedisLockLockFailedErr = errors.New("lock failed")
	// RedisLockTimeoutErr 加锁超时
	RedisLockTimeoutErr = errors.New("timeout")
)
