package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"grpc-todolist/config"
	"strings"
	"time"
)

// 全局唯一, 不对外暴露
var _db *gorm.DB

func InitDB() {
	host := config.Conf.Mysql.Host
	port := config.Conf.Mysql.Port
	database := config.Conf.Mysql.Database
	username := config.Conf.Mysql.Username
	password := config.Conf.Mysql.Password
	charset := config.Conf.Mysql.Charset
	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=",
		charset, "&parseTime=true"}, "")
	err := Database(dsn)
	if err != nil {
		panic(err)
	}
	migration()
}

func Database(dsn string) error {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,  // 禁用 datetime 的进度, 兼容配置
		DontSupportRenameColumn:   true,  // 重命名索引, 采用删除并新建的方式, 兼容配置
		DontSupportRenameIndex:    true,  // 重命名索引, 采用删除并新建的方式, 兼容配置
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,  // 对于表的命名策略, 单数
		},
	})
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  // 设置连接池
	sqlDB.SetMaxOpenConns(100)  // 设置最大连接数
	sqlDB.SetConnMaxLifetime(30*time.Second)  // 设置最长连接时
	_db = db
	return nil
}

// NewDBClient 对外暴露 db 连接, 使用 context 进行进程控制
func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
