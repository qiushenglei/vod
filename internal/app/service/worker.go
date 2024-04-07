package service

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"vod/internal/app/service/stream"
)

var Worker *ants.MultiPool

func InitWorker() {
	var err error
	Worker, err = ants.NewMultiPool(100, 5, ants.RoundRobin)
	if err != nil {
		panic(err)
	}

	stream.Sch = make(chan stream.SchMsg, 100)

	Worker.Submit(func() {
		fmt.Println("listen")
		for {
			select {
			case msg := <-stream.Sch:
				stream.Scheduler(msg)
			}
		}
	})
}
