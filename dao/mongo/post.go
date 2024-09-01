package mongo

import (
	"comyun/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func AdjustFields(fields,postId string) error {
	//存入数据
	update := bson.M{
		"$set":bson.M{
			"fields":fields,
		},
	}
	postID,err := primitive.ObjectIDFromHex(postId)
	if err != nil{
		zap.L().Error("postId ObjectIDFromHex error"+err.Error())
		return err
	}
	_,err = CollPost.UpdateByID(context.TODO(),postID,update)
	if err != nil{
		zap.L().Error("Insert post error"+err.Error())
		return err
	}

	return nil
}

func GetFields(id string)(*models.FormFieldGroup,error){
	objectID,err := primitive.ObjectIDFromHex(id)
	if err != nil{
		zap.L().Error("transform to objectID error,"+err.Error())
		return nil,err
	}

	fields := new(models.FormFieldGroup)
	err = CollPost.FindOne(context.TODO(),bson.M{"_id":objectID}).Decode(fields)
	if err != nil{
		zap.L().Error("find PostField error,"+err.Error())
		return nil,err
	}

	return fields,nil
}

func SaveUserFields(uf *models.UserFieldGroup) error {
	if _,err := CollPost.InsertOne(context.TODO(),*uf);err != nil{
		zap.L().Error("insert UserFieldGroup error,"+err.Error())
		return err
	}

	return nil
}

func GetUserPosts(pl *models.PostListParams) ([]models.UserFieldGroup,error){
	filter := bson.M{"post_id":pl.PostID}

	//设置查询条件
	findOptions := options.Find()
	findOptions.SetSkip(int64((pl.Page-1)*pl.Num))
	findOptions.SetLimit(int64(pl.Num))
	findOptions.SetSort(bson.M{"submission_date":-1})

	//查询对应变量
	ctx := context.Background()
	cursor,err := CollPost.Find(ctx,filter,findOptions)
	if err != nil{
		return nil,err
	}
	defer cursor.Close(ctx)

	uf := new([]models.UserFieldGroup)
	if err = cursor.All(ctx,uf);err != nil{
		return nil,err
	}

	return *uf,nil
}
