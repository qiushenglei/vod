package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/q191201771/lal/pkg/base"
	"vod/internal/app/entity"
	"vod/internal/app/global"
	"vod/internal/app/service/stream"
	"vod/internal/xlal"
)

func BindUserSession(c *gin.Context, req *entity.BindGroupReq) (*stream.Session, error) {
	// 获取用户信息
	userId := c.GetString("UserID")

	// 创建session group

	session, err := stream.GetSession(userId)
	if err == nil && session != nil {
		return session, nil
	}

	// 2. 配置session
	sessionContext, err := xlal.LalServer.AddCustomizePubSession(userId)
	if err != nil {
		return nil, err
	}
	sessionContext.WithOption(func(option *base.AvPacketStreamOption) {
		option.VideoFormat = base.AvPacketStreamVideoFormatAnnexb
	})

	session = stream.NewSession(stream.WithSessionContext(sessionContext), stream.WithChan())
	session.AddStream(userId)

	// 用户绑定group
	stream.BindMap[userId] = session

	return session, nil
}

func DoPlay(c *gin.Context, req *entity.DoPlayReq) error {

	// 绑定用户session
	session, err := BindUserSession(c, &entity.BindGroupReq{})
	if err != nil {
		return err
	}

	// 如果有播放内容，关闭播放内容，重新加载
	if session.IsPlaying() {
		session.CmdChan <- stream.NewCommand(0, stream.ReloadCommand)
	}

	msg := stream.SchMsg{
		Session: session,
		Command: stream.NewCommand(0, stream.PlayCommand),
	}

	stream.Sch <- msg

	return nil
}

func SendCommand(c *gin.Context, req *entity.SendCommandReq) error {
	session, err := stream.GetSession(c.GetString(global.UserIDKey))
	if err != nil {
		return err
	}

	cmd := &stream.Command{}
	if err := copier.Copy(cmd, req); err != nil {
		return err
	}

	session.CmdChan <- cmd
	return nil
}
