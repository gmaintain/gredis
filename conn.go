package gredis

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	DefaultBufferSize = 64
)

var (
	ErrBadTcpConn   = errors.New("无效的tcp连接")
	ErrResponseType = fmt.Errorf("返回结果类型无效")
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

func (c *RedisConn) Call(command string, args ...interface{}) (interface{}, error) {
	var err error
	if err = c.request(command, args); err != nil {
		return nil, err
	}
	if err = c.BuffioWriter.Flush(); err != nil {
		return nil, err
	}
	return nil, err
}

func (c *RedisConn) request(command string, args []interface{}) error {
	var err error
	if err = c.writeLen('*', 1+len(args)); err != nil {
		return err
	}
	if err = c.writeString(command); err != nil {
		return err
	}
	for _, arg := range args {
		switch data := arg.(type) {
		case string:
			err = c.writeString(data)
		case int64:
			err = c.writeInt64(data)
		case int:
			err = c.writeInt64(int64(data))
		case bool:
			if data == true {
				err = c.writeString("1")
			} else {
				err = c.writeString("0")
			}
		default:
			err = c.writeString(fmt.Sprintf("%v", data))
		}
	}
	return err
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

//协议第二步,写入命令长度以及命令
func (c *RedisConn) writeString(command string) error {
	var err error
	if err = c.writeLen('$', len(command)); err != nil {
		return err
	}
	if _, err = c.BuffioWriter.WriteString(command + "\r\n"); err != nil {
		return err
	}
	return nil
}

func (c *RedisConn) writeInt64(n int64) error {
	return c.writeBytes(strconv.AppendInt([]byte{}, n, 10))
}

func (c *RedisConn) writeBytes(b []byte) error {
	var err error
	if err = c.writeLen('$', len(b)); err != nil {
		return err
	}
	if _, err = c.BuffioWriter.Write(b); err != nil {
		return err
	}
	if _, err = c.BuffioWriter.WriteString("\r\n"); err != nil {
		return err
	}
	return nil
}
