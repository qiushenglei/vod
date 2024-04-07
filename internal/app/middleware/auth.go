package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiushenglei/gin-skeleton/pkg/errorpkg"
	"vod/internal/app/data"
	"vod/internal/app/global"
	"vod/internal/app/service"
	"vod/internal/app/utils"
)

func Auth() func(*gin.Context) {
	return func(c *gin.Context) {

		var token string
		var err error
		if token, err = c.Cookie("token"); err != nil {
			utils.Response(c, err.Error(), errorpkg.NewBizErrx(global.AuthFailCode, global.AuthFailMessage))
			c.Abort()
			return
		}

		// 解析token
		var parseArr []string
		if parseArr, err = service.ParseTokenStr(token); err != nil {
			utils.Response(c, err.Error(), errorpkg.NewBizErrx(global.AuthFailCode, global.AuthFailMessage))
			c.Abort()
			return
		}

		// 验证token
		if res, err := data.VODRedis.Get(c, parseArr[0]).Result(); err != nil || res == "" || res != parseArr[1] {
			utils.Response(c, global.AuthFailCode, errorpkg.NewBizErrx(global.AuthFailCode, global.AuthFailMessage))
			c.Abort()
			return
		} else {
			fmt.Println(res)
		}

		// todo::根据token查询用户信息，保存到ctx
		c.Set(global.UserIDKey, parseArr[0])
		fmt.Print(token)

		c.Next()
	}
}
