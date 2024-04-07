package controller

import (
	"github.com/gin-gonic/gin"
	"vod/internal/app/entity"
	"vod/internal/app/global"
	"vod/internal/app/service"
	"vod/internal/app/utils"
)

func BindUserSession(c *gin.Context) {
	body := &entity.BindGroupReq{}
	if err := c.ShouldBind(body); err != nil {
		utils.Response(c, nil, err)
		return
	}

	_, err := service.BindUserSession(c, body)

	utils.Response(c, global.StatusOkMessage, err)
}

func DoPlay(c *gin.Context) {
	body := &entity.DoPlayReq{}
	if err := c.ShouldBind(body); err != nil {
		utils.Response(c, nil, err)
		return
	}

	err := service.DoPlay(c, body)

	utils.Response(c, global.StatusOkMessage, err)
}

func SendCommand(c *gin.Context) {
	body := &entity.SendCommandReq{}
	if err := c.ShouldBind(body); err != nil {
		utils.Response(c, nil, err)
		return
	}

	err := service.SendCommand(c, body)

	utils.Response(c, global.StatusOkMessage, err)
}
