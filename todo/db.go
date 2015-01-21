package todo

import (
	"fmt"
	"labix.org/v2/mgo"
)

const (
	DATABASE_NAME = "cpbr8app"
)

var conn *MongoConnection

func init() {
	conn = new(MongoConnection)
}

func GetMongoConnection() *MongoConnection {
	return conn
}

func GetNewConnection() (*mgo.Database, *mgo.Session) {
	session := conn.session.Copy()
	return session.DB(DATABASE_NAME), session
}

type MongoConnection struct {
	session *mgo.Session
}

func (m *MongoConnection) Start(user string, password string) {
	uri := fmt.Sprintf("mongodb://%s:%s@ds031711.mongolab.com:31711/%s",
		user, password, DATABASE_NAME)

	session, err := mgo.Dial(uri)
	if err != nil {
		fmt.Println("Authentication failed!")
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	m.session = session
}

func (m *MongoConnection) Close() {
	m.session.Close()
}
