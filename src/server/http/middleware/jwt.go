package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	TokenExpired     error  = errors.New("token is expired")
	TokenNotValidYet error  = errors.New("token not active yet")
	TokenMalformed   error  = errors.New("that's not even a token")
	TokenInvalid     error  = errors.New("couldn't handle this token")
	SignKey          string = "sasa"
)

//JWTLoginAuth middleware, check login toke
func JWTLoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.Next()
			return
		}

		j := NewJWT()

		//parse information from token
		claims, err := j.ParserToken(token)

		if err != nil {
			//token is outdated
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": -1,
					"msg":    "token is outdated, pls reply for authentication",
					"data":   nil,
				})
				c.Abort()
				return
			}
			//other errors
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    err.Error(),
				"data":   nil,
			})
			c.Abort()
			return
		}

		//parse info of a certain claims
		c.Set("token", token)
		c.Set("claims", claims)
	}
}

//JWTAuth middleware, check modification token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    "your request doesn't contain token, you have no right to access our system",
				"data":   nil,
			})
			c.Abort()
			return
		}

		j := NewJWT()

		//parse infomation from token
		claims, err := j.ParserToken(token)

		if err != nil {
			//token is outdated
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": -1,
					"msg":    "token is outdated, pls reply for authentication",
					"data":   nil,
				})
				c.Abort()
				return
			}
			//other errors
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": -1,
				"msg":    err.Error(),
				"data":   nil,
			})
			c.Abort()
			return
		}

		//parse info of a certain claims
		c.Set("token", token)
		c.Set("claims", claims)
	}
}

//JWT basic data structure
type JWT struct {
	SigningKey []byte
}

//define loading burden
type LoginClaims struct {
	UserID   int64
	Username string
	jwt.StandardClaims
}

//initialize JWT instance
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

//get signkey
func GetSignKey() string {
	return SignKey
}

//generate token based on user's infomation
//use HS256 algorithm to generate token
//use user's basic info and signkey to generate token
func (j *JWT) GenerateToken(claims LoginClaims) (string, error) {

	//input userID, username, expire duration into token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//SecretKey used for user data signing, cannot be exposed
	return token.SignedString(j.SigningKey)
}

//parse token
//couldn't handle this token
func (j *JWT) ParserToken(tokenString string) (*LoginClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		//jwt.ValidationError is an invalid toke structure
		if ve, ok := err.(*jwt.ValidationError); ok {
			//that means token is usable
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
				//ValidationErrorExpired means token outdated
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
				//ValidationErrorNotValidYet means useless token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	//parse information from claims and authenticate it with user's info
	if claims, ok := token.Claims.(*LoginClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

//update token
func (j *JWT) UpdateToken(tokenString string) (string, error) {
	// TimeFunc is a current time variable with default value time.Now, which is used for parsing token and outdated
	// token checking
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	//get data from token
	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	//authenticate whether this token is valid or not
	if claims, ok := token.Claims.(*LoginClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		//modify outdated time of claims
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.GenerateToken(*claims)
	}
	return "", fmt.Errorf("token获取失败:%v", err)
}
