package httproute

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vod/internal/app/controller"
	"vod/internal/app/middleware"
	"vod/internal/app/utils"
)

func InitHttpRoute(engine *gin.Engine) {
	engine.GET("/health", func(c *gin.Context) {
		utils.Response(c, http.StatusOK, nil)
	})

	unlogin := engine.Group("auth")
	{
		unlogin.POST("login", controller.Login)
	}

	// 短视频
	short := engine.Group("short", middleware.Auth())
	{
		short.POST("bindUserSession", controller.BindUserSession)
		short.POST("doplay", controller.DoPlay)
		short.POST("sendCommand", controller.SendCommand)
	}

	// 点播
	vod := engine.Group("vod", middleware.Auth())
	{
		vod.POST("bindUserSession", controller.BindUserSession)
	}

}
