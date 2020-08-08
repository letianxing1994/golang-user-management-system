package connection

import (
	pb "Entry_Task/src/public/protos"
	"Entry_Task/src/server/http/common"
	"log"
)

type TcpBridge interface {
	//connect tcp server
	Conn() error
}

type TcpManager struct {
	addr          string
	ServiceClient *pb.UserServiceClient
}

func NewTcpManager(addr string, serviceClient *pb.UserServiceClient) TcpBridge {
	return &TcpManager{addr: addr, ServiceClient: serviceClient}
}

func (t *TcpManager) Conn() (err error) {
	if t.addr == "" {
		t.addr = "localhost:50055"
	}
	if t.ServiceClient == nil {
		conn, err := common.NewTcpConn(t.addr)
		if err != nil {
			log.Fatal(err)
		}
		//defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		t.ServiceClient = &client
	}
	return
}
