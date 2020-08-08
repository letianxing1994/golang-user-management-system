package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"sync"
)

var mutexDB sync.Mutex
var wg1 sync.WaitGroup

func writeToDb(filename string, db *sql.DB) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("file read failed")
		return
	}

	input := bufio.NewScanner(file)
	idx := 1
	for input.Scan() {
		//mutexDB.Lock()
		str := input.Text()
		usrPwd := strings.Split(str, " ")
		stmt, err := db.Prepare("insert into user_tab (username, password) values (?, ?)")
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err1 := stmt.Exec(usrPwd[0], usrPwd[1])
		if err1 != nil {
			log.Fatal("sql syntax wrong 2")
			return
		}
		if idx%200000 == 0 {
			fmt.Printf("%s\t%d\t%s\n", filename, idx, " records has been written into database")
		}
		idx++
		//mutexDB.Unlock()
	}

	//wg1.Done()
}

func main() {
	//wg1.Add(5)
	db, err := sql.Open("mysql", "root:811149@Tim@/user_db")
	if err != nil {
		log.Fatal(err)
		return
	}

	writeToDb("/Applications/golang_test/user_info/user4.txt", db)
	//for i := 0; i < 5; i++ {
	//	go writeToDb("/Applications/golang_test/user_info/user"+strconv.Itoa(i)+".txt", db)
	//}
	//wg1.Wait()
	db.Close()
	fmt.Printf("%s", "数据库写成功")
}
