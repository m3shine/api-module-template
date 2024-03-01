package consts

import (
	"time"
)

type RedisStoreKey string

const (
	RedisPrefix = "sqrd:"
	RedisStoreTtl = 30 * 60 * time.Second
	
)
