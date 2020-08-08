package router

import (
	. "Entry_Task/src/server/http/apis"
	md "Entry_Task/src/server/http/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(userService *UserService) *gin.Engine {
	router := gin.New()

	//log middlewares
	//router.Use(gin.Recovery())
	//router.Use(middleware.LoggerMiddleware())

	//load template
	router.LoadHTMLGlob("/Users/tianxingle/go/src/Entry_Task/src/server/http/web/view/*")
	router.Static("/web/css", "/Users/tianxingle/go/src/Entry_Task/src/server/http/web/css")
	router.Static("/web/js", "/Users/tianxingle/go/src/Entry_Task/src/server/http/web/js")

	//route each http request to certain api
	router.GET("/login", (*userService).IndexApi)
	router.POST("/login", md.JWTLoginAuth(), (*userService).LogInApi)
	router.POST("/management/index", (*userService).ManagementApi)
	router.POST("/management/modify", md.JWTAuth(), (*userService).ModifyNickNameApi)
	router.POST("/management/upload", md.JWTAuth(), (*userService).UploadProfileApi)

	return router
}
