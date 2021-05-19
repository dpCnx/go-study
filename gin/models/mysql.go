package models

import (
	"fmt"

	"go-study/gin/conf"
	"go-study/gin/utils"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

)

var (
	db  *gorm.DB
	err error
)

func init() {

	db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.C.MySQL.User,
		conf.C.MySQL.Password,
		conf.C.MySQL.IP+":"+conf.C.MySQL.Port,
		conf.C.MySQL.Database)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	_ = db.Use(&utils.OpentracingPlugin{})

	zap.L().Debug("mysql start")
}
