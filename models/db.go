package models

import (
	"database/sql"
	"fmt"
	"github.com/yunling101/ControllerManager/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var (
	DB    *gorm.DB
	SqlDB *sql.DB
	err   error
)

func Connect() error {
	DB, err = gorm.Open(
		mysql.Open(common.Config().Global.Dsn()), &gorm.Config{
			Logger:         logger.Default.LogMode(logger.Silent),
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
		})
	if err != nil {
		return fmt.Errorf("connect db fail: %s", err.Error())
	}

	SqlDB, err = DB.DB()
	if err != nil {
		return fmt.Errorf("open db fail: %s", err.Error())
	}
	if err = SqlDB.Ping(); err != nil {
		return fmt.Errorf("ping db fail: %s", err.Error())
	}

	SqlDB.SetConnMaxLifetime(100 * time.Millisecond)
	SqlDB.SetMaxIdleConns(100) // 设置连接池的空闲数大小
	SqlDB.SetMaxOpenConns(100) // 设置最大打开连接数

	return nil
}
