package common

import (
	"Entry_Task/src/public/config"
	"golang.org/x/exp/errors/fmt"
)

func readConfig() (map[string]string, error) {
	//build map for config arguments storing
	record := make(map[string]string)

	//read config from file
	initParser := config.IniParser{}
	confFile := "Entry_Task/src/config/database.ini"
	if err := initParser.Load(confFile); err != nil {
		fmt.Printf("try load config file[%s] error[%s]\n", confFile, err.Error())
		return record, err
	}

	//get mysql config
	mysqlDrive := initParser.GetString("mysql", "db_drive_name")
	mysqlDB := initParser.GetString("mysql", "db_list")
	mysqlHost := initParser.GetString("mysql", "db_host")
	mysqlPort := initParser.GetString("mysql", "db_port")
	mysqlUser := initParser.GetString("mysql", "db_user")
	mysqlPassword := initParser.GetString("mysql", "db_password")
	mysqlMaxOpenConn := initParser.GetString("mysql", "db_max_open_conn")
	mysqlMaxIdleConn := initParser.GetString("mysql", "db_max_idle_conn")
	mysqlConnMaxLifeTime := initParser.GetString("mysql", "db_conn_max_life_time")

	//get redis cache config
	redisHost := initParser.GetString("redis", "rdb_host")
	redisPort := initParser.GetString("redis", "rdb_port")
	redisDB := initParser.GetString("redis", "rdb_DB")
	redisPassword := initParser.GetString("redis", "rdb_password")
	redisPoolSize := initParser.GetString("redis", "rdb_PoolSize")

	//record mysql key-value pairs
	record["db_drive_name"] = mysqlDrive
	record["db_list"] = mysqlDB
	record["db_host"] = mysqlHost
	record["db_port"] = mysqlPort
	record["db_user"] = mysqlUser
	record["db_password"] = mysqlPassword
	record["db_max_open_conn"] = mysqlMaxOpenConn
	record["db_max_idle_conn"] = mysqlMaxIdleConn
	record["db_conn_max_life_time"] = mysqlConnMaxLifeTime

	//record redis key-value pairs
	record["rdb_host"] = redisHost
	record["rdb_port"] = redisPort
	record["rdb_DB"] = redisDB
	record["rdb_password"] = redisPassword
	record["rdb_PoolSize"] = redisPoolSize

	return record, nil
}
