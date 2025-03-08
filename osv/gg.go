package osv

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gsadism/open-admin/pkg/image"
	"strconv"
	"strings"
	"sync"
	"time"
)

type _Image struct {
	client *image.GG
	cache  *redis.Client
	once   sync.Once
}

var Image = new(_Image)

func (i *_Image) Init(client *image.GG, cache *redis.Client) {
	if client != nil && cache != nil {
		i.once.Do(func() {
			i.client = client
			i.cache = cache
		})
	}
}

func (i *_Image) Number(n int) (string, string, error) {
	if code, img, err := i.client.Default(image.Number, n); err != nil {
		return "", "", err
	} else {
		// 将值存入缓存
		key := fmt.Sprintf("image_%s", strconv.FormatInt(SnowFlake.NextID(), 10))
		if err := i.cache.Set(context.TODO(), key, strings.ToLower(code), 3*time.Minute).Err(); err != nil {
			return "", "", err
		} else {
			return key, img, nil
		}
	}
}

func (i *_Image) Verify(key string, code string) bool {
	if vcode, err := i.cache.Get(context.TODO(), key).Result(); err != nil {
		return false
	} else {
		if strings.ToLower(vcode) != code {
			return false
		} else {
			return true
		}
	}
}
