package data

import (
	"log"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xmysql"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/comeonjy/go-kit/pkg/xmongo"
	"spider/configs"
)

var ProviderSet = wire.NewSet( NewData, NewWorkRepo)

type Data struct {
	Mongo *mongo.Collection
}

func newAccountMysql(cfg configs.Interface, logger *xlog.Logger) *gorm.DB {
	db := xmysql.New(cfg.Get().MysqlConf, logger)
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		log.Fatalln("AutoMigrate AccountModel err:", err)
	}
	return db
}
func NewData(cfg configs.Interface) *Data {
	return &Data{
		Mongo: newMongo(cfg),
	}
}
