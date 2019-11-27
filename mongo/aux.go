package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func IsNotFound(err error) bool {
	return err == mgo.ErrNotFound
}

func notFound() error {
	return mgo.ErrNotFound
}

func IdQ(id string) bson.M {
	return bson.M { "_id": bson.ObjectIdHex(id) }
}

func IdSafeQ(q bson.M, id string) bson.M {
	return IdSafeQ2(q, id, "_id")
}

func IdSafeQ2(q bson.M, id, field string) bson.M {
	oid, ok := ObjectId(id)
	if !ok {
		return nil
	}

	if q == nil {
		q = bson.M{}
	}

	q[field] = oid
	return q
}

func ObjectId(oid string) (bson.ObjectId, bool) {
	if bson.IsObjectIdHex(oid) {
		return bson.ObjectIdHex(oid), true
	} else {
		return "", false
	}
}
