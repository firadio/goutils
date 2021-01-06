package utils

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetDB(Host string, DB string) *mgo.Database {
	//连接
	session, err := mgo.Dial(Host)
	if err != nil {
		return nil
	}
	//获取文档集
	db := session.DB(DB)
	return db
}

func EnsureIndex(db *mgo.Database, collection string, key []string, unique bool) {
	// 创建索引
	index := mgo.Index{
		Key:        key,    // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     unique, // 唯一索引 同mysql唯一索引
		DropDups:   false,  // 索引重复替换旧文档,Unique为true时失效
		Background: true,   // 后台创建索引
	}
	if err := db.C(collection).EnsureIndex(index); err != nil {
		fmt.Println(err)
		return
	}
}

func DbInsert(db *mgo.Database, collection string, data interface{}) {
	if err := db.C(collection).Insert(data); err != nil {
		fmt.Println(err)
		return
	}
}

func DbInsert2(coll *mgo.Collection, data interface{}) {
	if err := coll.Insert(data); err != nil {
		fmt.Println(err)
		return
	}
}

func DbUpdateById(db *mgo.Database, collection string, id bson.ObjectId, data bson.M) {
	selector := bson.M{"_id": id}
	data = bson.M{"$set": data}
	if err := db.C(collection).Update(selector, data); err != nil {
		fmt.Println(err)
		return
	}
}

func DbGetOneByInt(db *mgo.Database, collection string, fieldName string, searchValue int, OneRow interface{}) error {
	if err := db.C(collection).Find(bson.M{fieldName: searchValue}).One(OneRow); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
