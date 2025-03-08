package osv

import (
	"github.com/gsadism/open-admin/pkg/storage"
	"sync"
)

type _Minio struct {
	client *storage.Minio
	domain string
	once   sync.Once
}

var Minio = new(_Minio)

func (m *_Minio) Init(client *storage.Minio, domain string) {
	if client != nil {
		m.once.Do(func() {
			m.client = client
			m.domain = domain
		})
	}
}

func (m *_Minio) Client() *storage.Minio {
	return m.client
}

func (m *_Minio) Domain() string {
	return m.domain
}
