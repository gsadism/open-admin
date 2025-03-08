package osv

import (
	"github.com/gsadism/open-admin/pkg/next/snowflake"
	"sync"
)

type _SnowFlake struct {
	client *snowflake.Worker
	once   sync.Once
}

var SnowFlake = new(_SnowFlake)

func (s *_SnowFlake) Init(client *snowflake.Worker) {
	if client != nil {
		s.once.Do(func() {
			s.client = client
		})
	}
}

func (s *_SnowFlake) NextID() int64 {
	return s.client.GetID()
}
