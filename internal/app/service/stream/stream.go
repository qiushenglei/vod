package stream

import (
	"vod/internal/app/entity"
)

const (
	LocalSourceType int = iota
	RemoteSourceType
)

var BindMap map[string]*Session

func init() {
	BindMap = make(map[string]*Session)
}

type Stream interface {
	StreamList() []*entity.Stream
	Play([]*entity.Stream) error
}

func NewStream(sourceType int, userId string, session *Session) (Stream, error) {
	var s Stream
	var err error

	switch sourceType {
	case LocalSourceType:
		s, err = NewLocalStream(userId, session)
	case RemoteSourceType:
		//lalSession = &RemoteStream{}
	}
	return s, err
}
