package utils

import (
	"fmt"
	"github.com/q191201771/lal/pkg/aac"
	"github.com/q191201771/lal/pkg/avc"
	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/logic"
	"github.com/q191201771/naza/pkg/nazalog"
	"os"
	"time"
)

// ReadAudioPacketsFromFile 从aac es流文件读取所有音频包
func ReadAudioPacketsFromFile(filename string) (audioContent []byte, audioPackets []base.AvPacket) {
	var err error
	audioContent, err = os.ReadFile(filename)
	nazalog.Assert(nil, err)

	pos := 0
	timestamp := float32(0)
	for {
		ctx, err := aac.NewAdtsHeaderContext(audioContent[pos : pos+aac.AdtsHeaderLength])
		nazalog.Assert(nil, err)

		packet := base.AvPacket{
			PayloadType: base.AvPacketPtAac,
			Timestamp:   int64(timestamp),
			Payload:     audioContent[pos+aac.AdtsHeaderLength : pos+int(ctx.AdtsLength)],
		}

		audioPackets = append(audioPackets, packet)

		timestamp += float32(48000*4*2) / float32(8192*2) // (frequence * bytePerSample * channel) / (packetSize * channel)

		pos += int(ctx.AdtsLength)
		if pos == len(audioContent) {
			break
		}
	}

	return
}

// ReadVideoPacketsFromFile 从h264 es流文件读取所有视频包
func ReadVideoPacketsFromFile(filename string) (videoContent []byte, videoPackets []base.AvPacket) {
	var err error
	videoContent, err = os.ReadFile(filename)
	nazalog.Assert(nil, err)

	timestamp := float32(0)
	err = avc.IterateNaluAnnexb(videoContent, func(nal []byte) {
		// 将nal数据转换为lalserver要求的格式输入
		packet := base.AvPacket{
			PayloadType: base.AvPacketPtAvc,
			Timestamp:   int64(timestamp),
			Payload:     append(avc.NaluStartCode4, nal...),
		}

		videoPackets = append(videoPackets, packet)

		t := avc.ParseNaluType(nal[0])
		if t == avc.NaluTypeSps || t == avc.NaluTypePps || t == avc.NaluTypeSei {
			// noop
			fmt.Println(t, " eq")
		} else {
			fmt.Println(t)
			timestamp += float32(1000) / float32(15) // 1秒 / fps
		}
	})
	nazalog.Assert(nil, err)

	return
}

// mergePackets 将音频队列和视频队列按时间戳有序合并为一个队列
func mergePackets(audioPackets, videoPackets []base.AvPacket) (packets []base.AvPacket) {
	var i, j int
	for {
		// audio数组为空，将video的剩余数据取出，然后merge结束
		if i == len(audioPackets) {
			packets = append(packets, videoPackets[j:]...)
			break
		}

		//
		if j == len(videoPackets) {
			packets = append(packets, audioPackets[i:]...)
			break
		}

		// 音频和视频都有数据，取时间戳小的
		if audioPackets[i].Timestamp < videoPackets[j].Timestamp {
			packets = append(packets, audioPackets[i])
			i++
		} else {
			packets = append(packets, videoPackets[j])
			j++
		}
	}

	return
}

// normalPlay 按时间戳间隔匀速发送音频和视频
func normalPlay(packets []base.AvPacket, session logic.ICustomizePubSessionContext) {
	// 4. 按时间戳间隔匀速发送音频和视频
	startRealTime := time.Now()
	startTs := int64(0)
	for i := range packets {

		// 数据包的时间(元数据的时间)
		diffTs := packets[i].Timestamp - startTs

		// 当前时间 - 开始传输文件的时间 = 已传输的时长(如果播放了10s，那么这样就是10*1000)
		diffReal := time.Now().Sub(startRealTime).Milliseconds()
		//nazalog.Debugf("%d: %s, %d, %d", i, packets[i].DebugString(), diffTs, diffReal)

		// sleep等待前一个包内容播放完毕后，再发下一个包给到客户端
		if diffReal < diffTs {
			time.Sleep(time.Duration(diffTs-diffReal) * time.Millisecond)
		}
		session.FeedAvPacket(packets[i])
	}
}

// TimingPlay 定时播放
func TimingPlay(packets []base.AvPacket, session logic.ICustomizePubSessionContext, n int64) {
	// 4. 按时间戳间隔匀速发送音频和视频
	startRealTime := time.Now()
	startTs := n
	for i := range packets {

		//if packets[i].Timestamp <= n*1000 {
		//	continue
		//}

		// 数据包的时间(元数据的时间)
		diffTs := packets[i].Timestamp - startTs

		// 当前时间 - 开始传输文件的时间 = 已传输的时长(如果播放了10s，那么这样就是10*1000)
		diffReal := time.Now().Sub(startRealTime).Milliseconds()
		//nazalog.Debugf("%d: %s, %d, %d", i, packets[i].DebugString(), diffTs, diffReal)

		// sleep等待前一个包内容播放完毕后，再发下一个包给到客户端
		if diffReal < diffTs {
			time.Sleep(time.Duration(diffTs-diffReal) * time.Millisecond)
		}
		err := session.FeedAvPacket(packets[i])
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}
