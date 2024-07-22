package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"golang-pet-api/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func InitGorm() {
	var err error
	var dbConfig = global.Config.Db

	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Db,
		dbConfig.Charset,
		dbConfig.Timeout)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // 禁用彩色打印
		},
	)

	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	if global.Db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         newLogger,
	}); err != nil {
		zap.S().Panic("mysql连接失败, error=" + err.Error())
	}

	global.Db.Debug()
	
	// 连接成功
	zap.S().Debugf("musq连接成功")
	sqlDb, _ := global.Db.DB()
	sqlDb.SetMaxIdleConns(dbConfig.MaxIdle)
	sqlDb.SetMaxOpenConns(dbConfig.MaxOpen)

}
