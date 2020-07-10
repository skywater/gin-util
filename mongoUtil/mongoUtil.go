package mongoUtil

import (
	"context"
	"log"
	"time"

	"github.com/skywater/gin-util/stringUtil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 完整url=mongodb://root:123456@localhost:27017/
type mgo struct {
	url        string // 完整url = mongodb://root:123456@localhost:27017/
	URI        string // 数据库网络地址，localhost:27017/
	User       string // 账号
	Pwd        string // 密码
	database   string // 要连接的数据库
	collection string // 要连接的集合
	MgoColl    *mongo.Collection
}

func (m *mgo) Init() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //养成良好的习惯，在调用WithTimeout之后defer cancel()
	if stringUtil.IsBlank(m.url) {
		m.url = "mongodb://"
		if stringUtil.IsBlank(m.URI) {
			m.URI = "localhost:27017/"
		}
		if stringUtil.IsBlank(m.User) {
			m.url += m.URI
		} else {
			m.url += m.User + ":" + m.Pwd + "@" + m.URI
		}
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.url))
	if err != nil {
		log.Println(err)
		panic(err)
	}
	// defer client.Disconnect(ctx)
	return client
}

func (m *mgo) InitCollection() *mongo.Collection {
	if nil != m.MgoColl {
		database := m.Init().Database(m.database)
		collection := database.Collection(m.collection)
		m.MgoColl = collection
	}
	return m.MgoColl
}

func (m *mgo) MgoCollection(databaseName string, collectionName string) *mongo.Collection {
	if stringUtil.IsBlank(databaseName) {
		databaseName = m.database
	}
	if stringUtil.IsBlank(collectionName) {
		collectionName = m.collection
	}
	database := m.Init().Database(databaseName)
	collection := database.Collection(collectionName)
	m.MgoColl = collection
	return collection
}

func (m *mgo) Find(databaseName string, collectionName string) (*mongo.Collection, []map[string]interface{}) {
	collection := m.MgoCollection(databaseName, collectionName)
	filter := bson.D{}
	// SORT := bson.D{{"_id", 1}} //filter := bson.D{{key,value}}
	// findOptions := options.Find().SetSort(SORT).SetLimit(3).SetSkip(1)
	// findOptions := options.Find().SetLimit(3).SetSkip(1)
	cursor, e := collection.Find(context.TODO(), filter, nil)
	// cursor, e := collection.Find(context.TODO(), filter, findOptions)
	if e != nil {
		panic(e)
	}
	var episodes []map[string]interface{}
	// var episodes []Episode
	e = cursor.All(context.TODO(), &episodes)
	if e != nil {
		panic(e)
	}
	// log.Println(ToPrettyJson(episodes))
	return collection, episodes
}

func (m *mgo) InsertOne(v interface{}) (*mongo.Collection, interface{}) {
	Result, error := m.InitCollection().InsertOne(context.TODO(), v)
	if nil != error {
		log.Println(error)
	}
	return m.MgoColl, Result.InsertedID
}

func (m *mgo) InsertMany(v ...interface{}) *mongo.Collection {
	_, error := m.InitCollection().InsertMany(context.TODO(), v)
	if nil != error {
		log.Println(error)
	}
	return m.MgoColl
}

// Delete 方法，v为nil，也删除全部！！！
func (m *mgo) Delete(key string, v interface{}) *mongo.Collection {
	filter := bson.D{}
	if nil != v {
		filter = bson.D{{key, v}}
	}
	// filter不能为nil，nil不处理，删除全部用bson.D{}
	_, error := m.InitCollection().DeleteMany(context.TODO(), filter)
	if nil != error {
		log.Println(error)
	}
	return m.MgoColl
}

func (m *mgo) Drop() {
	error := m.InitCollection().Drop(context.TODO())
	if nil != error {
		log.Println(error)
	}
}