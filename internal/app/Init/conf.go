package Init

import (
	"fmt"
	"github.com/qiushenglei/gin-skeleton/pkg/conf"
	"github.com/qiushenglei/gin-skeleton/pkg/safe"
	"vod/internal/app/entity"
)

var Configure entity.Configure

func InitConf() {
	// lal

	// vod
	filePath := GetConfFileName(Env, HttpConfFile)
	if err := conf.InitConf(&Configure, filePath); err != nil {
		panic("load HttpConfFile failed")
	}
}

func GetConfFileName(env, fileName string) string {
	relativeFileName := fmt.Sprintf("%s/%s_%s", "internal/app/conf/", env, fileName)
	absolutePath := safe.Path(relativeFileName)
	return absolutePath
}
