package gredis

import (
	"fmt"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
	du, _ := time.ParseDuration("2s")
	conn, err := Dial("127.0.0.1:6379", du)
	defer conn.Close()
	if err != nil {
		t.Log(err)
	}

	RedisConn := NewRedisConn(conn)
	//RedisConn.writeLen('*', 5)
	setval, err := RedisConn.Set("mykey", "myval")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(setval)
}
