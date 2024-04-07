package xlal

import (
	"context"
	"github.com/anguloc/zet/pkg/safe"
	"github.com/q191201771/lal/pkg/logic"
	"vod/internal/app/Init"
)

var LalServer logic.ILalServer

func Start(c context.Context) {

	// 启动lal服务
	LalServer = logic.NewLalServer(func(option *logic.Option) {
		option.ConfFilename = Init.GetConfFileName(Init.Env, Init.LalConfFile)
	})

	// todo:: lal添加sessionHook,实现rtmp转其他协议
	safe.Go(c, func(ctx context.Context) {
		LalServer.RunLoop()
	}, nil)
}
