package model

import (
	"fmt"
	"github.com/hxx258456/pyramidel-chain-baas/internal/localconfig"
	mlogger "github.com/hxx258456/pyramidel-chain-baas/pkg/utils/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db = connect()

var mysqlServeLogger = mlogger.Lg.Named("mysql/gorm")

func connect() *gorm.DB {
	conf := genMysqlConfig()
	gormLogger := mlogger.NewGormLogger()
	db, err := gorm.Open(mysql.New(conf), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		mysqlServeLogger.Panic("Couldn't connect to database", zap.Error(err))
	}
	db.AutoMigrate(&Host{})
	db.AutoMigrate(&Organization{})
	db.AutoMigrate(&Peer{})
	db.AutoMigrate(&Orderer{})
	db.AutoMigrate(&Member{})
	return db
}

func genMysqlConfig() mysql.Config {
	var parserTime, loc string
	if localconfig.Defaultconfig.Mysql.Parsetime {
		parserTime = "True"
	} else {
		parserTime = "False"
	}
	if localconfig.Defaultconfig.Mysql.Loc {
		loc = "Local"
	} else {
		loc = "UTC"
	}
	format := "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s"
	dns := fmt.Sprintf(format, localconfig.Defaultconfig.Mysql.User, localconfig.Defaultconfig.Mysql.Password, localconfig.Defaultconfig.Mysql.Host, localconfig.Defaultconfig.Mysql.Port, localconfig.Defaultconfig.Mysql.DB, localconfig.Defaultconfig.Mysql.Charset, parserTime, loc)

	mysqlServeLogger.Info(dns)
	return mysql.Config{
		DSN:                       dns,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
}
