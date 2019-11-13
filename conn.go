package gredis

import (
	"bufio"
	"errors"
	"log"
	"net"
	"time"
)

const (
	DefaultBufferSize = 64
)

var (
	ErrBadTcpConn = errors.New("无效的tcp连接")
)

type RedisConn struct {
	TcpConn      *net.TCPConn
	Buffer       []byte
	BuffioWriter *bufio.Writer
}

func NewRedisConn(tcpConn *net.TCPConn) *RedisConn {
	return &RedisConn{
		TcpConn:      tcpConn,
		Buffer:       make([]byte, DefaultBufferSize),
		BuffioWriter: bufio.NewWriter(tcpConn),
	}
}

func Dial(address string, timeout time.Duration) (*net.TCPConn, error) {
	c, err := net.DialTimeout("tcp", address, timeout)
	conn, ok := c.(*net.TCPConn)
	if !ok {
		return nil, ErrBadTcpConn
	}
	if err != nil {
		log.Printf("net dial err %#v", err)
		return nil, err
	}
	return conn, nil
}

// 协议格式,写入命令或者参数的长度
func (c *RedisConn) writeLen(prefix byte, n int) error {
	pos := len(c.Buffer) - 1
	c.Buffer[pos] = '\n'
	pos--
	c.Buffer[pos] = '\r'
	pos--
	//这里为了提高性能,115 在[]byte中存3个byte,先计算几个位,再byte转换为每一位的ascii值
	for i := n; i != 0 && pos >= 0; i = i / 10 {
		c.Buffer[pos] = byte(i%10 + '0')
		pos--
	}
	c.Buffer[pos] = prefix
	_, err := c.BuffioWriter.Write(c.Buffer[pos:])
	return err
}
