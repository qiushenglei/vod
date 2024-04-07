package stream

import (
	"fmt"
	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/logic"
	"github.com/qiushenglei/gin-skeleton/pkg/errorpkg"
	"sync"
	"time"
	"vod/internal/app/Init"
	"vod/internal/app/global"
)

type Session struct {
	stream     Stream
	CmdChan    chan *Command
	lalSession logic.ICustomizePubSessionContext
	lock       *sync.Mutex
	playFlag   bool
}

type Command struct {
	Timestamp int64
	Sign      CommandType
}

type CommandType int

const (
	PlayCommand     CommandType = 0
	PauseCommand    CommandType = 1 // 暂停
	ContinueCommand CommandType = 2 // 继续
	MoveCommand     CommandType = 3 // 移动进度条
	ReloadCommand   CommandType = 4 // 重新加载
)

var (
	PlayingDuration int64 = 10 // 播放时长
	WaitDuration    int64 = 8  // 等待时长
)

type Option func(*Session)

func NewSession(op ...Option) *Session {
	s := &Session{}
	for _, v := range op {
		v(s)
	}

	s.lock = &sync.Mutex{}
	return s
}

func WithChan() Option {
	return func(s *Session) {
		s.CmdChan = make(chan *Command, 10)
	}
}

func WithSessionContext(ses logic.ICustomizePubSessionContext) Option {
	return func(s *Session) {
		s.lalSession = ses
	}
}

func (s *Session) AddStream(userId string) {
	st, _ := NewStream(Init.Configure.Source, userId, s)
	s.stream = st
}

func NewCommand(pos int64, sign CommandType) *Command {
	return &Command{
		pos,
		sign,
	}
}

func GetSession(userId string) (*Session, error) {
	if _, ok := BindMap[userId]; !ok {
		return nil, errorpkg.NewBizErrx(global.SessionFailCode, global.SessionFailMessage)
	}
	return BindMap[userId], nil
}

func (s *Session) IsPlaying() bool {
	return s.playFlag
}

func (s *Session) DoPlay() {

	// 获取视频策略
	list := s.stream.StreamList()

	// 播放，推流
	s.stream.Play(list)

}

// play 定时播放
func play(packets []base.AvPacket, session *Session, step int64) {
	// 4. 按时间戳间隔匀速发送音频和视频
	//startRealTime := time.Now()
	//startTs := n
	//tick := time.NewTimer(time.Duration(duration) * time.Millisecond)

	var multiple int64

	packageCount := len(packets)

	fmt.Sprintf("开始倒流")
	for i := 0; i < packageCount; i++ {
		// 如果接收到客户端指令，要发送对应的数据包
		select {
		case command := <-session.CmdChan:
			_, multiple, i = session.actionCommand(command, packets, i)
		default:

		}

		// 每发送N秒后，sleepM秒。
		tmpMul := packets[i].Timestamp / (PlayingDuration * 1000)
		if tmpMul > multiple {
			fmt.Sprintf("开始休息:%d秒 idx:%d timestamp:%d \n", WaitDuration, i, packets[i].Timestamp)
			time.Sleep(time.Duration(WaitDuration) * time.Second)
			fmt.Sprintf("结束休息:%d秒 idx:%d timestamp:%d \n", WaitDuration, i, packets[i].Timestamp)
			multiple = tmpMul
		}

		err := session.lalSession.FeedAvPacket(packets[i])
		if err != nil {
			fmt.Print(err.Error())
		}
	}

}

func (s *Session) actionCommand(command *Command, packets []base.AvPacket, nowPackageIdx int) (CommandType, int64, int) {
	var mul int64
	var packageIdx int
	action := command.Sign

	fmt.Println("exec Command ", action, nowPackageIdx)
	switch action {
	case PauseCommand:
		<-s.CmdChan
		packageIdx = nowPackageIdx - 1
		fmt.Println("pause end, continue")
	case ContinueCommand:
		mul = command.Timestamp / (PlayingDuration * 1000)
		var item base.AvPacket
		for packageIdx, item = range packets {
			if item.Timestamp <= command.Timestamp {
				break
			}
		}
	case MoveCommand:
		mul = command.Timestamp / (PlayingDuration * 1000)
		var item base.AvPacket
		for packageIdx, item = range packets {
			if item.Timestamp <= command.Timestamp {
				break
			}
		}
	case ReloadCommand:
		// 直接跳到最后一个包，结束
		packageIdx = len(packets)
		mul = packets[packageIdx-1].Timestamp / (PlayingDuration * 1000)
	}

	// -1是为了回到上一个包，避免不连续
	return action, mul, packageIdx - 1
}
