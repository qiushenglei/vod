package stream

import (
	"bytes"
	"fmt"
	"github.com/q191201771/lal/pkg/base"
	"github.com/qiushenglei/gin-skeleton/pkg/safe"
	"github.com/yapingcat/gomedia/go-mp4"
	"io"
	"os"
	"vod/internal/app/entity"
)

type LocalStream struct {
	Session *Session
}

var CodeMap = map[mp4.MP4_CODEC_TYPE]base.AvPacketPt{
	mp4.MP4_CODEC_H264:  base.AvPacketPtAvc,
	mp4.MP4_CODEC_H265:  base.AvPacketPtAvc,
	mp4.MP4_CODEC_AAC:   base.AvPacketPtAac,
	mp4.MP4_CODEC_G711A: base.AvPacketPtG711A,
	mp4.MP4_CODEC_G711U: base.AvPacketPtG711U,
}

func NewLocalStream(userId string, session *Session) (*LocalStream, error) {
	return &LocalStream{
		Session: session,
	}, nil
}

func (s *LocalStream) StreamList() []*entity.Stream {
	return []*entity.Stream{
		{
			Name: "1",
			Url:  safe.Path("data/video/4.mp4"),
		},
		//{
		//	Name: "2",
		//	Url:  safe.Path("data/video/2.mp4"),
		//},
	}
}

func (s *LocalStream) Play(list []*entity.Stream) error {
	s.Session.playFlag = true
	defer func() {
		s.Session.playFlag = false
		fmt.Println("playing end")
	}()
	for _, item := range list {
		// _, nalu := utils.ReadVideoPacketsFromFile(item.Url)
		//utils.TimingPlay(nalu, s.Session.lalSession, 0)
		//lals.DelCustomizePubSession(session)

		//
		f, err := os.Open(item.Url)
		if err != nil {
			fmt.Println(err)
			return err
		}

		demuxer := mp4.CreateMp4Demuxer(f)
		if infos, err := demuxer.ReadHead(); err != nil && err != io.EOF {
			fmt.Println(err)
		} else {
			fmt.Printf("%+v\n", infos)
		}

		mp4info := demuxer.GetMp4Info()
		fmt.Printf("%+v\n", mp4info)

		var flvByte []byte
		flvBuff := bytes.NewBuffer(flvByte)
		//flvWrite := flv.CreateFlvWriter(flvBuff)

		var pkgList []base.AvPacket
		for {
			pkg, err := demuxer.ReadPacket()
			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Printf("track:%d,cid:%+v,pts:%d dts:%d\n", pkg.TrackId, pkg.Cid, pkg.Pts, pkg.Dts)
			if pkg.Cid == mp4.MP4_CODEC_H264 {
				//flvWrite.WriteH264(pkg.Data, uint32(pkg.Pts), uint32(pkg.Dts))
				//vfile.Write(pkg.Data)
				fmt.Println(265)
			} else if pkg.Cid == mp4.MP4_CODEC_H265 {
				fmt.Println(265)
			} else if pkg.Cid == mp4.MP4_CODEC_AAC {
				fmt.Println("AAC")
				//afile.Write(pkg.Data)
			} else if pkg.Cid == mp4.MP4_CODEC_MP3 {
				fmt.Println("MP3")
				//afile.Write(pkg.Data)
			}

			if _, ok := CodeMap[pkg.Cid]; !ok {
				panic("unknown codec")
			}

			basePkg := base.AvPacket{
				PayloadType: CodeMap[pkg.Cid],
				Timestamp:   int64(pkg.Dts),
				Pts:         int64(pkg.Pts),
				Payload:     pkg.Data,
			}
			flvBuff.Write(pkg.Data)

			//s.Session.lalSession.FeedAvPacket(basePkg)
			pkgList = append(pkgList, basePkg)
		}

		play(pkgList, s.Session, 0)

	}

	return nil
}
