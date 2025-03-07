package osv

import (
	"context"
	"gorm.io/gorm"
	"sync"
)

type db struct {
	client *gorm.DB
	once   sync.Once
}

var DB = new(db)

func (d *db) Init(client *gorm.DB) {
	if client != nil {
		d.once.Do(func() {
			d.client = client
		})
	}
}

func (d *db) WithContext(ctx context.Context) *gorm.DB {
	return d.client.Set("context", ctx)
}
