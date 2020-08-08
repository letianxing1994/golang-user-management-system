package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func write2000Tofile(src string, dst string) {
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
	idx := 1
	for input.Scan() && idx <= 2000 {
		str := input.Text()
		usrPwd := strings.Split(str, " ")
		tmp = usrPwd[0] + "," + "123456" + "\n"
		needed := len(tmp)
		written, err3 := file2.WriteString(tmp)
		if needed != written || err3 != nil {
			log.Fatal("current user and password info is writtern failed")
			break
		}
		idx++
	}

	file1.Close()
	file2.Close()
}

func main() {
	write2000Tofile("/Applications/golang_test/user_info/user0.txt", "/Applications/golang_test/user_info/user_test3.csv")
	fmt.Printf("%s", "finished")
}
