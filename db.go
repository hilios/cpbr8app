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

// Returns the current global connection object
func GetMongoConnection() *MongoConnection {
	return conn
}

// Returns a the project database and the current session
func GetDatabase() (*mgo.Database, *mgo.Session) {
	s := conn.session.Copy()
	return s.DB(DATABASE), s
}

// Stores the global mongo session
type MongoConnection struct {
	session *mgo.Session
}

// Connect to the dabase with given user and password and store the session
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

// Close the stored session
func (m *MongoConnection) Close() {
	m.session.Close()
}
