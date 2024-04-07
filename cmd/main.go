package main

import (
	"context"
	"fmt"
	"github.com/anguloc/zet/pkg/safe"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	Init "vod/internal/app/Init"
	"vod/internal/app/data"
	"vod/internal/app/service"
	"vod/internal/xlal"
	httproute "vod/router/http"
)

// ffmpeg -i "./data/test.flv" -r 25  -f flv "rtmp://127.0.0.1:1935/live/hello"
// http://127.0.0.1:8080/live/hello.flv
// ffplay rtmp://127.0.0.1/live/hello

func main() {
	c, _ := context.WithCancel(context.Background())

	// flag
	Init.InitFlag()

	// conf
	Init.InitConf()

	// db
	data.InitData(c)

	// worker
	service.InitWorker()

	// lal server
	xlal.Start(c)

	// vod server
	engine := gin.Default()
	httproute.InitHttpRoute(engine)

	httpListener, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", Init.Configure.Domain, Init.Configure.Port))

	safe.Go(c, func(ctx context.Context) {
		http.Serve(httpListener, engine)
	}, nil)

	// ListenSignal 监听信号

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	fmt.Println("挂起服务启动协程") // graceful shutdown
	<-quit
	fmt.Println("\nShutdown all server ...")

}
