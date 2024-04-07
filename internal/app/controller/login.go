package controller

import (
	"github.com/gin-gonic/gin"
	"vod/internal/app/entity"
	"vod/internal/app/service"
	"vod/internal/app/utils"
)

func Login(c *gin.Context) {
	body := &entity.LoginReq{}
	if err := c.ShouldBind(body); err != nil {
		utils.Response(c, nil, err)
		return
	}

	data, err := service.Login(c, body)

	utils.Response(c, data, err)

}
