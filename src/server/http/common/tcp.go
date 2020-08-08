package common

import (
	"google.golang.org/grpc"
)

//build tcp connection
func NewTcpConn(addr string) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(addr, grpc.WithInsecure())
	return
}
