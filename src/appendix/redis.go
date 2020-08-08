package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/go-redis/redis"
)

var ctx = context.Background()

func ExampleNewClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)
}

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ok, err := rdb.HMSet(ctx, "0",
		"user_id", "1",
		"username", "letianxing",
		"password", "123456",
		"nickname", "sasa",
		"profile_picture", "www.sasa.com").Result()
	if err != nil && !ok {
		fmt.Println(err.Error() + "  错误1")
	}

	val, err := rdb.HMGet(ctx,
		"0", "user_id", "nickname", "profile_picture", "password", "username").Result()
	if err != nil {
		fmt.Println(err.Error() + "  错误2")
	}

	for _, v := range val {
		fmt.Println(v.(string))
	}

	//val2, err := rdb.Get(ctx, "key2").Result()
	//if err == redis.Nil {
	//	fmt.Println("key2 does not exist")
	//} else if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("key2", val2)
	//}
}

func main() {
	var usrPwd []rune
	var cryptoUserPwd [32]byte
	userPassword := ""

	for i := 1; i < 9; i++ {
		usrPwd = append(usrPwd, rune(i))
		if i == 8 {
			strPwd := string(usrPwd)
			cryptoUserPwd = sha256.Sum256([]byte(strPwd))
			userPassword = fmt.Sprintf("%x", cryptoUserPwd)
			fmt.Printf("%s\n", userPassword)
			fmt.Printf("%s", strPwd)
		}
	}
}
