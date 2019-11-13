package gredis

import (
	"bufio"
	"log"
	"net"
)

func Get(conn *net.TCPConn, str string) (string, error) {

	reader := bufio.NewReader(conn)
	_, _ = conn.Write([]byte(str + "\r\n"))

	//strResp,_ := reader.ReadString('\n')
	strResp, _ := reader.ReadBytes('\n')
	log.Printf("read string resp: %#v", string(strResp))
	return string(strResp), nil
}
