package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Jresp struct {
	data   data   `json:"data"`
	msg    string `json:"msg"`
	status int    `json:"status"`
}

type data struct {
	token          string `json:"token"`
	nickname       string `json:"nickname"`
	profilePicture string `json:"profilePicture"`
}

func main() {
	file, err := os.Open("/Applications/golang_test/user_info/user_test3.csv")
	if err != nil {
		log.Fatal(err)
		return
	}
	file2, err1 := os.OpenFile("/Applications/golang_test/user_info/token_2000_test.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err1 != nil {
		log.Fatal(err1)
		return
	}
	//var j Jresp

	input := bufio.NewScanner(file)
	for input.Scan() {
		client := &http.Client{}
		str := input.Text()
		usrPwd := strings.Split(str, ",")
		req, err := http.NewRequest("POST",
			"http://localhost:8080/login",
			strings.NewReader(url.Values{"username": {usrPwd[0]}, "password": {usrPwd[1]}}.Encode()))
		if err != nil {
			fmt.Println(err)
			return
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("read body err, %v\n", err)
			return
		}

		//errj := json.Unmarshal([]byte(body), &j)
		strRes := string([]byte(body))
		strs := strings.Split(strRes, ":")
		tokenNick := strs[2]
		token := strings.Split(tokenNick, ",")[0]

		file2.WriteString(token[1:len(token)-1] + "," + "suiji" + "\n")
		resp.Body.Close()
	}
	file.Close()
	file2.Close()
}
