package osv

import (
	"github.com/go-redis/redis/v8"
	"sync"
)

type _Redis struct {
	client *redis.Client
	once   sync.Once
}

var Redis = new(_Redis)

func (r *_Redis) Init(client *redis.Client) {
	if client != nil {
		r.once.Do(func() {
			r.client = client
		})
	}
}

func (r *_Redis) Client() *redis.Client {
	return r.client
}
