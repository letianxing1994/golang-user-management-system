package repositories

import (
	"Entry_Task/src/server/tcp/common"
	"Entry_Task/src/server/tcp/datamodels"
	"crypto/sha256"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	//connect database
	Conn() error
	UpdateNickName(userId int64, nickname string) error
	UpdateProfile(userId int64, profilePicture string) error
	SelectByUser(username, password string) (*datamodels.User, error)
}

type UserManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewUserManager(table string, db *sql.DB) UserRepository {
	return &UserManager{table: table, mysqlConn: db}
}

//mysql connection
func (u *UserManager) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user_tab"
	}
	return
}

//user nickname update
func (u *UserManager) UpdateNickName(userId int64, nickname string) error {
	//1. judge whether connection exists
	if err := u.Conn(); err != nil {
		return err
	}

	query := "update " + u.table + " set nickname = ? where user_id = ?"

	stmt, err := u.mysqlConn.Prepare(query)
	if stmt == nil {
		return nil
	}
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(nickname, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserManager) UpdateProfile(userId int64, profilePicture string) error {
	//1. judge whether connection exists
	if err := u.Conn(); err != nil {
		return err
	}

	query := "update " + u.table + " set profile_picture = ? where user_id = ?"

	stmt, err := u.mysqlConn.Prepare(query)
	if stmt == nil {
		return nil
	}
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(profilePicture, userId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserManager) SelectByUser(username, password string) (userRes *datamodels.User, err error) {
	//judge whether connection exists
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	query := "select * from " + u.table + " where username = ? and password = ?"
	row, errRow := u.mysqlConn.Query(query, username, fmt.Sprintf("%x", sha256.Sum256([]byte(password))))
	if row == nil {
		return &datamodels.User{}, nil
	}
	defer row.Close()
	if errRow != nil {
		return &datamodels.User{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, nil
	}
	userRes = &datamodels.User{}
	common.DataToStructByTagSql(result, userRes)
	return
}
