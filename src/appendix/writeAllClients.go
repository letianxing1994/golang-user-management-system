package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func writeTofile(src string, dst string) {
	file1, err1 := os.Open(src)
	if err1 != nil {
		log.Fatal("src file is broken")
		return
	}
	file2, err2 := os.OpenFile(dst, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err2 != nil {
		log.Fatal("dst file is broken")
		return
	}

	input := bufio.NewScanner(file1)
	tmp := ""
	for input.Scan() {
		str := input.Text()
		usrPwd := strings.Split(str, " ")
		tmp = usrPwd[0] + "," + "123456" + "\n"
		needed := len(tmp)
		written, err3 := file2.WriteString(tmp)
		if needed != written || err3 != nil {
			log.Fatal("current user and password info is writtern failed")
			break
		}
	}

	file1.Close()
	file2.Close()
}

func main() {
	for i := 0; i < 5; i++ {
		writeTofile("/Applications/golang_test/user_info/user"+strconv.Itoa(i)+".txt", "/Applications/golang_test/user_info/user_test2.csv")
		fmt.Printf("%d\t%s\n", i, " file has been written into database")
	}
}
