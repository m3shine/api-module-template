package model

import (
	"west.garden/template/common/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var db *gorm.DB

func Init(cfg *config.MysqlConfig) error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	idb, _ := db.DB()
	idb.SetConnMaxIdleTime(120 * time.Second)
	idb.SetConnMaxLifetime(7200 * time.Second)
	idb.SetMaxOpenConns(200)
	idb.SetMaxIdleConns(10)
	if err := idb.Ping(); err != nil {
		return err
	}
	return nil
}

func Close() error {
	if db != nil {
		idb, err := db.DB()
		if err == nil {
			idb.Close()
		}
	}
	fmt.Println("close mysql connect")
	return nil
}

func GetDB() *gorm.DB {
	return db
}
