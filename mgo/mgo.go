package mgo

import (
	"fmt"

	"github.com/firadio/goutils/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongo struct {
	Db *mgo.Database
}

func New(Host string, DB string) Mongo {
	mongo := Mongo{}
	//连接
	session, err := mgo.Dial(Host)
	if err != nil {
		return mongo
	}
	//获取文档集
	mongo.Db = session.DB(DB)
	return mongo
}

func (mongo Mongo) GetQueue(collection string, bsonM bson.M, iLimit int) ([]bson.M, error) {
	AllRows := []bson.M{}
	c := mongo.Db.C(collection)
	if err := c.Find(bsonM).Sort("queued").Limit(iLimit).All(&AllRows); err != nil {
		fmt.Println(err)
		return AllRows, err
	}
	arr := []interface{}{}
	for _, row := range AllRows {
		arr = append(arr, row["_id"])
	}
	selector := bson.M{"_id": bson.M{"$in": arr}}
	data := bson.M{"$set": bson.M{"queued": utils.TimestampInt32()}}
	_, err := c.UpdateAll(selector, data)
	if err != nil {
		fmt.Println(err)
		return AllRows, err
	}
	return AllRows, nil
}
