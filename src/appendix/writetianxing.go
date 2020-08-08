package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root:811149@Tim@/user_db")
	if err != nil {
		log.Fatal(err)
		return
	}

	stmt, err := db.Prepare("insert into user_tab (username, password, nickname) values (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err1 := stmt.Exec("letianxing", fmt.Sprintf("%x", sha256.Sum256([]byte("123456"))), "sasa")
	if err1 != nil {
		log.Fatal("sql syntax wrong 2")
		return
	}

	db.Close()
}
