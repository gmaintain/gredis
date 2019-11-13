package gredis

func (c *RedisConn) Get(key string) ([]byte, error) {
	_, err := c.Call("GET", key)
	if err != nil {
		return nil, err
	}
	return nil, err
}
