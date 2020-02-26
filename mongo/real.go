package mongo

import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/unectio/util"
)

type MgoSession struct {
	s *mgo.Session
}

func (ms *MgoSession)Close() {
	ms.s.Close()
}

func (ms *MgoSession)Copy() Session {
	return &MgoSession{ms.s.Copy()}
}

func (ms *MgoSession)Collection(loc *Location) Collection {
	return &MgoCollection{ms.s.DB(loc.Db).C(loc.Col)}
}

type MgoCollection struct {
	c *mgo.Collection
}

func (mc *MgoCollection)Find(q bson.M) Query {
	return &MgoQuery{mc.c.Find(q)}
}

func (mc *MgoCollection)Update(q bson.M, u interface{}) error {
	return mc.c.Update(q, u)
}

func (mc *MgoCollection)RemoveId(id bson.ObjectId) error {
	return mc.c.Remove(bson.M{"_id": id})
}

func (mc *MgoCollection)RemoveAll(q bson.M) error {
	_, err := mc.c.RemoveAll(q)
	if err == mgo.ErrNotFound {
		err = nil
	}
	return err
}

func (mc *MgoCollection)EnsureIndex(idx *mgo.Index) error {
	return mc.c.EnsureIndex(*idx)
}

func (mc *MgoCollection)Insert(q interface{}) error {
	return mc.c.Insert(q)
}

type MgoQuery struct {
	q *mgo.Query
}

func (mq *MgoQuery)Count() (int, error) {
	return mq.q.Count()
}

func (mq *MgoQuery)One(out interface{}) error {
	return mq.q.One(out)
}

func (mq *MgoQuery)Iter() Iter {
	return mq.q.Iter()
}

func (mq *MgoQuery)Sort(f string) Query {
	mq.q.Sort(f)
	return mq
}

func (mq *MgoQuery)Limit(n int) Query {
	mq.q.Limit(n)
	return mq
}

func Connect(url string) (Session, error) {
	c := util.CredsParse(url)

	info := mgo.DialInfo{
		Addrs:		[]string{c.Adr + ":" + c.Prt},
		Database:	c.Dom,
		Username:	c.Usr,
		Password:	c.Pwd,
	}

	log.Printf("-> [%s]\n", c.String())
	s, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}

	return &MgoSession{s}, nil
}
