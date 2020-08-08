package common

import (
	"context"
	"github.com/go-redis/redis"
)

//build redis connection
func NewRdsConn() (rdb *redis.Client, err error) {

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 250,
	})

	_, err = rdb.Ping(context.Background()).Result()
	return rdb, err
}

//get results, return one result
func GetCacheResultRow(val []interface{}) map[string]string {
	scanArgs := make([]interface{}, len(val))
	values := make([][]byte, len(val))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for i, v := range values {
		if v != nil {
			if i == 0 {
				record["user_id"] = string(v)
			} else if i == 1 {
				record["nickname"] = string(v)
			} else if i == 2 {
				record["profile_picture"] = string(v)
			} else if i == 3 {
				record["password"] = string(v)
			} else {
				record["username"] = string(v)
			}
		}
	}
	return record
}
