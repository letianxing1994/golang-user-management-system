package main

import (
	"Entry_Task/src/server/http/apis"
	"Entry_Task/src/server/http/connection"
	. "Entry_Task/src/server/http/router"
	"log"
)

func main() {
	//1.build tcp connection to tcp server
	tcpManager := connection.NewTcpManager("", nil)
	err := tcpManager.Conn()
	if err != nil {
		log.Fatal(err)
		return
	}
	userService := apis.NewUserServiceManager(tcpManager)

	//2.bind router
	router := InitRouter(&userService)

	err = router.Run()
	if err != nil {
		log.Fatalf("failed to open HTTP server: %v", err)
		return
	}
}
