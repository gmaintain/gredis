package gredis

func (c *RedisConn) Get(key string) ([]byte, error) {
	val, err := c.Call("GET", key)
	if err != nil {
		return nil, err
	}
	_, ok := val.([]byte)
	if !ok {
		return nil, ErrResponseType
	}
	return val.([]byte), err
}

func (c *RedisConn) Set(key, val string) ([]byte, error) {
	result, err := c.Call("SET", key, val)
	if err != nil {
		return nil, err
	}
	_, ok := result.([]byte)
	if !ok {
		return nil, ErrResponseType
	}
	return result.([]byte), err
}
