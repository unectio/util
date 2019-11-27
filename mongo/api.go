package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

type Location struct {
	Db	string
	Col	string
}

type Session interface {
	Collection(*Location) Collection
	Copy() Session
	Close()
}

type Collection interface {
	Find(q bson.M) Query
	Update(q bson.M, u interface{}) error
	RemoveId(id bson.ObjectId) error
	RemoveAll(q bson.M) error
	Insert(q interface{}) error
}

type Query interface {
	One(out interface{}) error
	Iter() Iter
	Sort(f string) Query
	Limit(n int) Query
	Count() (int, error)
}

type Iter interface {
	Next(out interface{}) bool
	Close() error
	Err() error
}
