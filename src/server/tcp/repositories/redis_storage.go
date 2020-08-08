package repositories

import (
	"Entry_Task/src/server/tcp/common"
	"Entry_Task/src/server/tcp/datamodels"
	"context"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type RedisStorage interface {
	//connect redis
	Conn() error
	GetByToken(token string) (*datamodels.User, error)
	UpdateNickNameByToken(token, nickname string) error
	UploadProfileByToken(token, profilePicture string) error
	SetByUser(token, username, password, nickname, profilePicture string, userID int64) error
}

type RedisManager struct {
	redisClient *redis.Client
}

func NewRedisManager(redisClient *redis.Client) RedisStorage {
	return &RedisManager{redisClient: redisClient}
}

//redis connection
func (r *RedisManager) Conn() (err error) {
	if r.redisClient == nil {
		rdb, err := common.NewRdsConn()
		if err != nil {
			return err
		}
		r.redisClient = rdb
	}
	return
}

//get user info from redis
func (r *RedisManager) GetByToken(token string) (userRes *datamodels.User, err error) {
	//judge whether connection exists
	if err = r.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	client := r.redisClient
	val, err := client.HMGet(context.Background(),
		token, "user_id", "nickname", "profile_picture", "password", "username").Result()
	if err != nil {
		return &datamodels.User{}, err
	}
	result := common.GetCacheResultRow(val)
	if len(result) == 0 {
		return &datamodels.User{}, err
	}
	userRes = &datamodels.User{}
	common.DataToStructByTagSql(result, userRes)
	return
}

//delete user from redis
func (r *RedisManager) UpdateNickNameByToken(token, nickname string) (err error) {
	//judge whether connection exists
	if err = r.Conn(); err != nil {
		return err
	}

	client := r.redisClient
	ok, err := client.HMSet(context.Background(), token, "nickname", nickname).Result()

	if err != nil && !ok {
		return err
	}
	return
}

func (r *RedisManager) UploadProfileByToken(token, profilePicture string) (err error) {
	//judge whether connection exists
	if err = r.Conn(); err != nil {
		return err
	}

	client := r.redisClient
	ok, err := client.HMSet(context.Background(), token, "profile_picture", profilePicture).Result()

	if err != nil && !ok {
		return err
	}
	return
}

//set user into redis
func (r *RedisManager) SetByUser(token, username, password, nickname, profilePicture string, userID int64) (err error) {
	//judge whether connection exists
	if err = r.Conn(); err != nil {
		return err
	}

	client := r.redisClient
	ok, err := client.HMSet(context.Background(),
		token, "user_id", strconv.FormatInt(userID, 10),
		"username", username,
		"password", password,
		"nickname", nickname,
		"profile_picture", profilePicture).Result()
	client.Expire(context.Background(), token, time.Duration(time.Now().Unix()+3600))

	if err != nil && !ok {
		return err
	}
	return
}
