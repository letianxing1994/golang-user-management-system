package main

import (
	pb "Entry_Task/src/public/protos"
	"Entry_Task/src/server/tcp/common"
	"Entry_Task/src/server/tcp/repositories"
	"Entry_Task/src/server/tcp/service"
	_ "Entry_Task/src/server/tcp/service"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50055"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		fmt.Println("TCP server starts, now it is accepting requests......")
	}

	//connect database
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Fatal(err)
		return
	}

	//connect redis
	rdb, err := common.NewRdsConn()
	if err != nil {
		log.Fatal(err)
		return
	}

	//register repository and grpc apis
	rpcServer := grpc.NewServer()
	userRepository := repositories.NewUserManager("user_tab", db)
	redisStorage := repositories.NewRedisManager(rdb)
	userService := service.NewUserServiceManager(userRepository, redisStorage)
	pb.RegisterUserServiceServer(rpcServer, userService)

	//grpc server starts
	if err := rpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
		return
	}
}
