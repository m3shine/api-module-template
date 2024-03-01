package consts

import (
	"errors"
	"time"
)

type LockKey string

const (
	//RedisLockTtl expire time
	RedisLockTtl = time.Second * 30
	//RedisLockResetTTLInterval redis lock reset ttl interval
	RedisLockResetTTLInterval = RedisLockTtl / 3
	//RedisLockTryLockInterval retry lock interval
	RedisLockTryLockInterval = time.Second
	//RedisLockUnlockScript lock script
	RedisLockUnlockScript = `
			if redis.call("get",KEYS[1]) == ARGV[1] then
				return redis.call("del",KEYS[1])
			else
				return 0
			end`

	//RedisKeyCPUAssemble cpu assemble key
	RedisKeyCPUAssemble LockKey = "sqrd:cpu:assemble:%s"
	RedisKeyGPUAssemble LockKey = "sqrd:gpu:assemble:%s"
	RedisKeyOpenBox     LockKey = "sqrd:open:box:%s"
	RedisKeyClaimInitRewards LockKey = "sqrd:claim:init:rewards:%s"
)

var (
	// RedisLockLockFailedErr lock failed
	RedisLockLockFailedErr = errors.New("lock failed")
	// RedisLockTimeoutErr lock timeout
	RedisLockTimeoutErr = errors.New("timeout")
)
