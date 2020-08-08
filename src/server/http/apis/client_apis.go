package apis

import (
	pb "Entry_Task/src/public/protos"
	"Entry_Task/src/server/http/connection"
	"Entry_Task/src/server/http/datamodels"
	md "Entry_Task/src/server/http/middleware"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

type UserService interface {
	IndexApi(c *gin.Context)
	LogInApi(c *gin.Context)
	ManagementApi(c *gin.Context)
	ModifyNickNameApi(c *gin.Context)
	UploadProfileApi(c *gin.Context)
}

type UserServiceManager struct {
	manager connection.TcpBridge
}

func NewUserServiceManager(manager connection.TcpBridge) UserService {
	return &UserServiceManager{manager: manager}
}

func LogInCheck(username, replyName, password, replyPwd string) bool {
	return username == replyName && fmt.Sprintf("%x", sha256.Sum256([]byte(password))) == replyPwd
}

func generateToken(c *gin.Context, reply *pb.LogInReply) (*datamodels.LoginResult, string) {
	//generate signkey
	j := md.NewJWT()

	//get user's info
	userID := reply.GetUserId()
	username := reply.GetUsername()
	nickname := reply.GetNickname()
	profilePicture := reply.GetProfilePicture()

	//build claims info
	claims := md.LoginClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000, //sign key valid begin time
			ExpiresAt: time.Now().Unix() + 3600, //sign key invalid time
			Issuer:    md.GetSignKey(),
		},
	}

	//generate token according to claims
	token, err := j.GenerateToken(claims)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    err.Error(),
			"data":   nil,
		})
		return nil, ""
	}

	//get user's info
	data := datamodels.LoginResult{
		Token:          token,
		Nickname:       nickname,
		ProfilePicture: profilePicture,
	}

	return &data, token
}

func (u *UserServiceManager) IndexApi(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (u *UserServiceManager) LogInApi(c *gin.Context) {
	bridge := u.manager
	manager := *bridge.(*connection.TcpManager)
	client := *manager.ServiceClient

	//check whether token exists
	token, ok := c.Get("token")

	//this request has no token
	if !ok {
		//authenticate user's info
		username := c.PostForm("username")
		password := c.PostForm("password")
		reply, err := client.LogIn(context.Background(),
			&pb.LogInRequest{
				Username: username,
				Password: password})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -1,
				"msg":    err.Error(),
				"data":   nil,
			})
			return
		}

		//check log in info
		isPass := LogInCheck(username, reply.GetUsername(), password, reply.GetPassword())
		if isPass {
			data, token := generateToken(c, reply)
			if token == "" {
				return
			}
			//write user into redis
			reply, err := client.MysqlToCache(context.Background(),
				&pb.MysqlToCacheRequest{
					Token:          token,
					UserId:         reply.GetUserId(),
					Username:       reply.GetUsername(),
					Password:       reply.GetPassword(),
					Nickname:       reply.GetNickname(),
					ProfilePicture: reply.GetProfilePicture(),
				})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status": -1,
					"msg":    err.Error(),
					"data":   nil,
				})
				return
			}
			//build json to front end
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"msg":    "log in successfully, " + reply.GetMessage(),
				"data":   &data,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    "validation failed, your username or password is incorrect",
				"data":   nil,
			})
		}
	} else {
		//use token to read info from redis
		reply, err := client.LogInCache(context.Background(),
			&pb.LogInCacheRequest{
				Token: token.(string)})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -1,
				"msg":    err.Error(),
				"data":   nil,
			})
			return
		}

		//get info from grpc reply
		nickname := reply.GetNickname()
		profilePicture := reply.GetProfilePicture()

		//if redis doesn't have this token, that means this token is outdated
		if err == redis.Nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    "token is outdated, pls reply for authentication",
				"data":   nil,
			})
		} else {
			//else get user's info from redis
			data := datamodels.LoginResult{
				Token:          token.(string),
				Nickname:       nickname,
				ProfilePicture: profilePicture,
			}

			//sending previous info
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"msg":    "this user has already logged in",
				"data":   data,
			})
		}
		return
	}
}

func (u *UserServiceManager) ManagementApi(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", nil)
}

func (u *UserServiceManager) ModifyNickNameApi(c *gin.Context) {
	bridge := u.manager
	manager := *bridge.(*connection.TcpManager)
	client := *manager.ServiceClient

	//get user's id from claims and update in mysql
	claim := c.MustGet("claims").(*md.LoginClaims)
	UserID := claim.UserID
	if c.PostForm("nickname") == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status": -1,
			"msg":    "nickname cannot be modified to empty",
		})
		return
	}
	reply, err := client.ModifyNickName(context.Background(),
		&pb.ModifyNickNameRequest{UserId: UserID, Nickname: c.PostForm("nickname")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
	} else {
		//update redis
		token := c.MustGet("token").(string)
		resp, err := client.ModifyNickNameCache(context.Background(),
			&pb.ModifyNickNameCacheRequest{Token: token, Nickname: c.PostForm("nickname")})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"msg":    resp.GetMessage() + " " + reply.GetMessage(),
			})
		}
	}
}

func (u *UserServiceManager) UploadProfileApi(c *gin.Context) {
	bridge := u.manager
	manager := *bridge.(*connection.TcpManager)
	client := *manager.ServiceClient

	//update mysql
	claim := c.MustGet("claims").(*md.LoginClaims)
	UserID := claim.UserID
	if c.PostForm("profilePicture") == "" {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status": -1,
			"msg":    "you cannot upload empty profile picture",
		})
		return
	}
	reply, err := client.UploadProfile(context.Background(),
		&pb.UploadProfileRequest{UserId: UserID, ProfilePicture: c.PostForm("profilePicture")})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
	} else {
		//update redis
		token := c.MustGet("token").(string)
		resp, err := client.UploadProfileCache(context.Background(),
			&pb.UploadProfileCacheRequest{Token: token, ProfilePicture: c.PostForm("profilePicture")})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": 0,
				"msg":    resp.GetMessage() + " " + reply.GetMessage(),
			})
		}
	}
}
