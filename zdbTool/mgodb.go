package zdbTool

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

func err_handler(err error) {
	fmt.Printf("err_handler, error:%s\n", err.Error())
	panic(err.Error())
}

func init() {
	mgodbConnect()
}

func mgodbConnect() *mgo.Session {

	// 初始化mongo dail info    后期配置到全局文件中

	dial_info := &mgo.DialInfo{
		Addrs:     []string{"127.0.0.1"},
		Direct:    false,
		Timeout:   time.Second * 1,
		Database:  "zinxmongo",
		Source:    "admin",
		Username:  "shv",
		Password:  "123456",
		PoolLimit: 1024,
	}

	session, err := mgo.DialWithInfo(dial_info)
	if err != nil {
		fmt.Printf("mgo dail error[%s]\n", err.Error())
		err_handler(err)
	}

	session.SetMode(mgo.Monotonic, true)

	defer session.Close()

	return session
}

// 获取文档数据  每一次操作都copy一份Session， 避免每次创建导致链接数量超过最大值
func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	ms := mgodbConnect().Copy()
	c := ms.DB(db).C(collection)
	/*
		setMode
		0 Eventual
		1 Monotonic
		2 Strong //default


		Strong 一致性模式
		session 的读写操作总向 primary 服务器发起并使用一个唯一的连接，因此所有的读写操作完全的一致（不存在乱序或者获取到旧数据的问题）。

		Monotonic 一致性模式
		session 的读操作开始是向某个 secondary 服务器发起（且通过一个唯一的连接），只要出现了一次写操作，session 的连接就会切换至 primary 服务器。
		由此可见此模式下，能够分散一些读操作到 secondary 服务器，但是读操作不一定能够获得最新的数据

		Eventual 一致性模式
		session 的读操作会向任意的 secondary 服务器发起，多次读操作并不一定使用相同的连接，也就是读操作不一定有序。
		session 的写操作总是向 primary 服务器发起，但是可能使用不同的连接，也就是写操作也不一定有序。Eventual 一致性模式最快，其是一种资源友好（resource-friendly）的模式。

	*/
	ms.SetMode(mgo.Monotonic, true)
	return ms, c
}

func Insert(db, collection string, doc interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Insert(doc)
}

func FindOne(db, collection, sort string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).Sort(sort).One(result)
}

func FindAll(db, collection, sort string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).Sort(sort).All(result)
}

func Update(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Update(selector, update)
}

// 更新，如果不存在就插入一个新的数据 `upsert:true`
func Upsert(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.Upsert(selector, update)
	return err
}

// `multi:true`
func UpdateAll(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.UpdateAll(selector, update)
	return err
}

func Remove(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Remove(selector)
}

func RemoveAll(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.RemoveAll(selector)
	return err
}

func FindPage(db, collection string, page, limit int, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).Skip(page * limit).Limit(limit).All(result)
}

func IsEmpty(db, collection string) bool {
	ms, c := connect(db, collection)
	defer ms.Close()
	count, err := c.Count()
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

func Count(db, collection string, query interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}
