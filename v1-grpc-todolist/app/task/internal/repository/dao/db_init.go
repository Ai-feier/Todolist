package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"grpc-todolist/config"
	logger2 "grpc-todolist/pkg/utils/logger"
	"strings"
	"time"
)

var _db *gorm.DB

func InitDB() {
	myConfig := config.Conf.Mysql
	host := myConfig.Host
	port := myConfig.Port
	database := myConfig.Database
	username := myConfig.Username
	password := myConfig.Password
	charset := myConfig.Charset
	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=", charset, "&parseTime=true"}, "")
	err := DataBase(dsn)
	if err != nil {
		logger2.LogrusObj.Error(err)
		panic(err)
	}
}

func DataBase(dsn string) error {
	// 初始化日志对象
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                           dsn,
		SkipInitializeWithVersion:     false, // 根据版本自动配置
		DefaultStringSize:             256,
		DisableDatetimePrecision:      true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:        true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:       true,
	}), &gorm.Config{
		Logger: ormLogger,  // 把自定义日志放进去
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,  // 单数命名策略
		},  
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)  // 
	sqlDB.SetMaxIdleConns(20) // 设置连接池，空闲
	sqlDB.SetConnMaxLifetime(30*time.Second)  // 一次连接的最大存货时长
	_db = db
	migration() 
	return err
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}















