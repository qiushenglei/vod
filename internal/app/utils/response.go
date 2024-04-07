package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/qiushenglei/gin-skeleton/pkg/errorpkg"
	"net/http"
	"vod/internal/app/entity"
	"vod/internal/app/global"
)

func Response(c *gin.Context, data interface{}, err error) {
	if err == nil {
		c.JSON(
			http.StatusOK,
			entity.DefaultResponse{
				Msg:  global.StatusOkMessage,
				Data: data,
			},
		)
	} else {
		// err放到context,熔断需要统计
		c.Error(err)

		// 获取error的code 和 msg
		var code int
		var msg string
		if e, ok := err.(errorpkg.Errx); ok {
			code = e.Code()
			msg = e.Msg()
		} else {
			code = errorpkg.CodeFalse
			msg = err.Error()
		}
		c.JSON(
			http.StatusInternalServerError,
			entity.DefaultResponse{
				Code: code,
				Msg:  msg,
				Data: data,
			},
		)
	}
}
