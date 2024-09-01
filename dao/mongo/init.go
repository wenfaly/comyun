package mongo

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

var ClientDB *mongo.Client
var CollPost *mongo.Collection

func Init() (error,context.Context) {
	//设置上下文
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d",
		viper.GetString("mongo.host"),
		viper.GetInt("mongo.port")),
	)

	client,err := mongo.Connect(ctx,clientOptions)
	if err != nil{
		zap.L().Error("Error to connect mongo "+err.Error())
		return err,ctx
	}

	err = client.Ping(ctx,nil)
	if err != nil{
		zap.L().Error("ping() error "+err.Error())
		return err,ctx
	}

	ClientDB = client
	CollPost = ClientDB.Database("post").Collection("coll_post")
	return nil,ctx
}

func Close(ctx context.Context){
	if err := ClientDB.Disconnect(ctx);err != nil{
		zap.L().Error("Disconnect mongoDB error "+err.Error())
		return
	}
}