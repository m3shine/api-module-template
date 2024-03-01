package redis


func Get(key string) ([]byte, error) {
	c := RedisClient.Get(ctx, key)
	return c.Bytes()
}
