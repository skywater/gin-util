package mongoUtil

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/skywater/gin-util/objectUtil"
	"github.com/skywater/gin-util/stringUtil"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MgoConfig 参数设置，完整url=mongodb://root:123456@localhost:27017/
type MgoConfig struct {
	URL      string // 完整url = mongodb://root:123456@localhost:27017/
	URI      string // 数据库网络地址，localhost:27017/
	User     string // 账号
	Pwd      string // 密码
	DbName   string // 要连接的数据库
	CollName string // 要连接的集合
	MgoColl  *mongo.Collection
	client   *mongo.Client
}

// Init 初始化链接
func (m *MgoConfig) Init() *mongo.Client {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if nil != m.client {
		return m.client
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //养成良好的习惯，在调用WithTimeout之后defer cancel()
	if stringUtil.IsBlank(m.URL) {
		m.URL = "mongodb://"
		if stringUtil.IsBlank(m.URI) {
			m.URI = "localhost:27017/"
		}
		if stringUtil.IsBlank(m.User) {
			m.URL += m.URI
		} else {
			m.URL += m.User + ":" + m.Pwd + "@" + m.URI
		}
	}
	log.Printf("mongodb 初始化链接：", m.URL)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.URL))
	if err != nil {
		log.Println(err)
		panic(err)
	}
	m.client = client
	// defer client.Disconnect(ctx)
	return client
}

// InitCollection 初始化集合
func (m *MgoConfig) InitCollection() *mongo.Collection {
	if nil == m.MgoColl {
		database := m.Init().Database(m.DbName)
		collection := database.Collection(m.CollName)
		m.MgoColl = collection
	}
	return m.MgoColl
}

// MgoCollection 初始化集合
func (m *MgoConfig) MgoCollection(databaseName string, collectionName string) *mongo.Collection {
	if stringUtil.IsBlank(databaseName) {
		databaseName = m.DbName
	}
	if stringUtil.IsBlank(collectionName) {
		collectionName = m.CollName
	}
	database := m.Init().Database(databaseName)
	collection := database.Collection(collectionName)
	m.MgoColl = collection
	return collection
}

// FindAll 查询所有
func (m *MgoConfig) FindAll() (*mongo.Collection, []map[string]interface{}) {
	collection := m.InitCollection()
	filter := bson.D{}
	cursor, e := collection.Find(context.TODO(), filter, nil)
	var retMap []map[string]interface{}
	e = cursor.All(context.TODO(), &retMap)
	if e != nil {
		panic(e)
	}
	return collection, retMap
}

// Find 查找集合中数据
func (m *MgoConfig) Find(databaseName string, collectionName string) (*mongo.Collection, []map[string]interface{}) {
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

// InsertMany 向集合中插入 map 数据，但不能直接接收 []map
func (m *MgoConfig) insertMany(v ...interface{}) *mongo.Collection {
	log.Println("InsertMany插入数据为：", v)
	if objectUtil.ArrayIsEmpty(v) {
		return m.MgoColl
	}
	_, error := m.InitCollection().InsertMany(context.TODO(), v)
	if nil != error {
		log.Println("insertMany插入数据异常：", error)
	}
	return m.MgoColl
}

// Insert 向集合中插入 map 数据，包括[]map
func (m *MgoConfig) Insert(v ...interface{}) *mongo.Collection {
	nv := objectUtil.DealArrayObj(v...)
	return m.insertMany(nv...)
}

// InsertJSON 向集合中插入 json 数据
func (m *MgoConfig) InsertJSON(v string) *mongo.Collection {
	log.Println("InsertJSON插入数据为：", v)
	if stringUtil.IsBlank(v) {
		log.Println("InsertJSON插入数据为空！")
		return m.MgoColl
	}
	v = strings.TrimSpace(v)
	if v[0:1] != "[" {
		v = "[" + v + "]"
	}
	var bDoc []interface{}
	bson.UnmarshalExtJSON([]byte(v), true, &bDoc)
	return m.insertMany(bDoc)
}

// DeleteAll 删除全部！！！
func (m *MgoConfig) DeleteAll() *mongo.Collection {
	return m.Delete("", nil)
}

// Delete 方法，v为nil时删除全部！！！
func (m *MgoConfig) Delete(key string, v interface{}) *mongo.Collection {
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

// Drop 删除表
func (m *MgoConfig) Drop() {
	error := m.InitCollection().Drop(context.TODO())
	if nil != error {
		log.Println(error)
	}
}

// Disconnect 断开链接
func (m *MgoConfig) Disconnect() {
	if nil != m.client {
		m.client.Disconnect(context.TODO())
		m.client = nil
	}
}
