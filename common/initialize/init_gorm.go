package initialize

import (
	"fmt"
	"golang-pet-api/common/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	global.Db, err = gorm.Open(mysql.Open(url), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	global.Db = global.Db.Debug()
	if err != nil {
		panic("mysql连接失败, error=" + err.Error())
	}
	if global.Db.Error != nil {
		panic(global.Db.Error)
	}
	// 连接成功
	fmt.Println("mysql连接成功!")
	sqlDb, err := global.Db.DB()
	sqlDb.SetMaxIdleConns(dbConfig.MaxIdle)
	sqlDb.SetMaxOpenConns(dbConfig.MaxOpen)

}
