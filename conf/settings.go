package conf

import (
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/middleware"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var MIDDLEWARE = []gin.HandlerFunc{
	gin.Recovery(),
	gin.Logger(),

	middleware.Cors(),
}

// GORM 配置
var GORM = &gorm.Config{
	// 是否跳过默认事务处理, 用于提高性能
	SkipDefaultTransaction: true,

	Logger: logger.Default.LogMode(logger.Info),
	// 命名策略
	NamingStrategy: schema.NamingStrategy{
		TablePrefix:   "",   // 所有表名都会有"t_"前缀
		SingularTable: true, // 使用单数表名
	},
	NowFunc: func() time.Time {
		return time.Now().Local() // 使用本地时间
	},
	AllowGlobalUpdate:                        false, // 允许全局更新
	PrepareStmt:                              false, // 为每个SQL语句创建一个prepared statement并缓存,提高执行效率
	DisableAutomaticPing:                     false, // 禁用自动ping
	DisableForeignKeyConstraintWhenMigrating: false, // 迁移时禁用外键约束
	DryRun:                                   false, // 开启调试模式，不会真正执行SQL
}
