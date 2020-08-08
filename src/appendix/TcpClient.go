package main

import (
	pb "Entry_Task/src/public/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	reply, err := client.UploadProfile(context.Background(),
		&pb.UploadProfileRequest{UserId: 10000004, ProfilePicture: "sasasa"})
	if err != nil {
		log.Fatal(err)
	}
	//strconv.FormatInt(reply.GetUserId()
	fmt.Printf("%s", reply.GetMessage())
}
