package data

import (
	"log"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xmysql"
	"github.com/google/wire"
	"gorm.io/gorm"

	"spider/configs"
)

var ProviderSet = wire.NewSet(NewData, NewTaskRepo, NewFetchRecordRepo, NewResourceRepo)

type Data struct {
	db *gorm.DB
}

func newSpiderMysql(cfg configs.Interface, logger *xlog.Logger) *gorm.DB {
	db := xmysql.New(cfg.Get().MysqlConf, logger)
	if err := db.AutoMigrate(&TaskModel{}); err != nil {
		log.Fatalln("AutoMigrate AccountModel err:", err)
	}
	if err := db.AutoMigrate(&FetchRecordModel{}); err != nil {
		log.Fatalln("AutoMigrate FetchRecordModel err:", err)
	}
	if err := db.AutoMigrate(&ResourceModel{}); err != nil {
		log.Fatalln("AutoMigrate ResourceModel err:", err)
	}
	return db
}

func NewData(cfg configs.Interface, logger *xlog.Logger) *Data {
	return &Data{
		db: newSpiderMysql(cfg, logger),
	}
}
