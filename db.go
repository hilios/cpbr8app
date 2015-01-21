package main

import (
	"fmt"
	"labix.org/v2/mgo"
)

const (
	DATABASE = "cpbr8app"
)

var conn *MongoConnection

func init() {
	conn = new(MongoConnection)
}

func GetMongoConnection() *MongoConnection {
	return conn
}

func GetDatabase() (*mgo.Database, *mgo.Session) {
	s := conn.session.Copy()
	return s.DB(DATABASE), s
}

type MongoConnection struct {
	session *mgo.Session
}

func (m *MongoConnection) Connect(user string, password string) {
	uri := fmt.Sprintf("mongodb://%s:%s@ds031711.mongolab.com:31711/%s",
		user, password, DATABASE)

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
