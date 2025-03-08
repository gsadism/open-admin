package conf

import (
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/core/base"
	"github.com/gsadism/open-admin/core/model"
	"github.com/gsadism/open-admin/logging"
	"github.com/gsadism/open-admin/middleware"
	"github.com/gsadism/open-admin/pkg/array"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

const (
	// SECRET_KEY : 系统秘钥,请妥善保存务泄露. 可通过go run main.go secret-generate-key 生成新的秘钥粘贴至此处
	SECRET_KEY = "b2440c871401ef5c82b44732fe714a18f8ffef6922eef9e7c4fb4e34d639a1c2"
)

var ROUTERS = []func(*gin.RouterGroup){
	base.Router,
}

var HOOKS = array.Merge[func()](
	base.Hooks,
)

var MODELS = array.Merge[model.IModel](
	base.Models,
)

var MIDDLEWARE = []gin.HandlerFunc{
	gin.Recovery(),
	gin.Logger(),

	middleware.Cors(),
}

/*====================== Logger ======================*/

// LogConsoleLevel : 控制台日志输出级别
var LogConsoleLevel = logging.DebugLevel

// LogFileLevel : 文件日志输出级别
var LogFileLevel = logging.ErrorLevel

// FileName : 日志文件名称
var FileName string = "open-admin.log"

// MaxSize : 日志文件达到的最大大小,以MB为单位
var MaxSize int = 1024

// MaxBackups : 留存的旧日志文件最大数量
var MaxBackups int = 30

// MaxAge : 旧日志文件的最大保存天数
var MaxAge int = 7

// Compress : 是否压缩日志文件
var Compress bool = false

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
