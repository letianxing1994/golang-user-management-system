package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

//var userName, userPassword string = "", ""
var info = make(map[string]string, 10000000)
var mutex sync.Mutex
var wg sync.WaitGroup

func getUsernameAndPassword(info *map[string]string) (bool, string, string) {
	usernameChar := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var eUsername []rune
	var cryptoUserPwd [32]byte
	userName := ""
	userPassword := ""

	//randomly generate 10-bits username
	for i := 0; i < 10; i++ {
		eUsername = append(eUsername, rune(usernameChar[rand.Intn(len(usernameChar))]))
		if i == 9 {
			userName = string(eUsername)
			if (*info)[userName] != "" {
				return false, "", ""
			}
		}
	}

	//randomly generate 6-bits password
	cryptoUserPwd = sha256.Sum256([]byte("123456"))
	userPassword = fmt.Sprintf("%x", cryptoUserPwd)
	(*info)[userName] = userPassword

	return true, userName, userPassword
}

func writeFile(info *map[string]string, filename string) {
	file, err :=
		os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	tmp := ""
	for i := 1; i <= 2000000; i++ {
		mutex.Lock()
		noneDup, u, p := getUsernameAndPassword(info)
		if !noneDup {
			i = i - 1
		} else {
			tmp = u + " " + p + "\n"
			needed := len(tmp)
			written, err1 := file.WriteString(tmp)
			if needed != written || err1 != nil {
				log.Fatal("current user and password info is written failed")
			} else {
				if i%200000 == 0 {
					fmt.Printf("%s\t%s\t%d\t%s\n", filename, " record", i, "is written")
				}
			}
		}
		mutex.Unlock()
	}

	file.Close()
	wg.Done()
}

func main() {
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go writeFile(&info, "/Applications/golang_test/user_info/user"+strconv.Itoa(i)+".txt")
	}

	wg.Wait()
	fmt.Printf("%s", "write succeed!!")
}
