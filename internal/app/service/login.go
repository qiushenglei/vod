package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiushenglei/gin-skeleton/pkg/errorpkg"
	"io"
	"strings"
	"time"
	"vod/internal/app/Init"
	"vod/internal/app/data"
	"vod/internal/app/entity"
	"vod/internal/app/global"
	"vod/internal/app/utils"
)

func Login(c context.Context, req *entity.LoginReq) (string, error) {
	// 对比,获取到用户uuid,这里先写死了
	userId := "hello"
	redisVal := utils.GenerateUniqueNumberByRand()
	TokenStr := GenerateTokenStr(fmt.Sprintf("%s_%d", userId, redisVal))

	// 保存token
	if r, err := data.VODRedis.Set(c, userId, redisVal, time.Hour*24).Result(); err != nil {
		return "", errorpkg.NewBizErrx(global.AuthFailCode, global.AuthFailMessage)
	} else {
		fmt.Println(r)
	}

	// 设置cookie
	if tc, ok := c.(*gin.Context); ok == true {
		tc.SetCookie("token", TokenStr, 3600*24, "/", Init.Configure.Domain, false, false)
	}

	return TokenStr, nil
}

func GenerateTokenStr(str string) string {
	// 方法1，可以替换base64的Encoding解码器
	buf := &bytes.Buffer{}
	writer := base64.NewEncoder(base64.StdEncoding, buf)
	if _, err := writer.Write([]byte(str)); err != nil {
		return ""
	}
	writer.Close()
	tokenVal := buf.String()

	// 方法2，使用标准base64编码器
	//tokenVal := base64.StdEncoding.EncodeToString([]byte(str))
	return tokenVal
}

func ParseTokenStr(base64Str string) ([]string, error) {

	var token []byte
	var err error

	// 方法1，可以替换base64的Encoding解码器
	buf := bytes.NewBufferString(base64Str)
	reader := base64.NewDecoder(base64.StdEncoding, buf)

	if token, err = io.ReadAll(reader); err != nil {
		return nil, err
	}

	// 方法2，使用标准base64编码器
	//token, err = base64.StdEncoding.DecodeString(base64Str)
	//if err != nil {
	//	return ""
	//}

	arr := strings.Split(string(token), "_")

	if len(arr) != 2 {
		return nil, errorpkg.NewBizErrx(global.AuthFailCode, global.AuthFailMessage)
	}
	return arr, nil
}
