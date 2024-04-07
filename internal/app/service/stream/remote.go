package stream

import "vod/internal/app/entity"

type RemoteStream struct {
}

func (s *RemoteStream) StreamList() []*entity.Stream {
	return nil
}
