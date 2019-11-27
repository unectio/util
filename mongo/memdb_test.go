package mongo

import (
	"fmt"
	"testing"
	"gopkg.in/mgo.v2/bson"
)

type Foo struct {
	Id	bson.ObjectId	`bson:"_id"`
	A	int		`bson:"a"`
	B	string		`bson:"b"`
}

func TestMemdb(t *testing.T) {
	mdb := GetMemDB()

	col := mdb.Collection(&Location{"foo", "bar"})
	err := col.Insert(&Foo{Id: bson.NewObjectId(), A: 1, B: "one"})
	if err != nil {
		fmt.Printf("error insert: %s\n", err.Error())
		t.FailNow()
	}

	col.Insert(&Foo{Id: bson.NewObjectId(), A: 2, B: "two"})

	var x Foo

	fmt.Printf("--- LOOKUP ---\n")
	err = col.Find(bson.M{"a": 1}).One(&x)
	if err != nil {
		fmt.Printf("error looking up: %s\n", err.Error())
		t.Fail()
	}

	fmt.Printf("A:%d B:%s (id: %v)\n", x.A, x.B, x.Id)
	fmt.Printf("--- DELETE ---\n")
	err = col.RemoveId(x.Id)
	if err != nil {
		fmt.Printf("error remove: %s\n", err.Error())
		t.Fail()
	}

	fmt.Printf("--- ITER ---\n")
	iter := col.Find(bson.M{}).Iter()
	for iter.Next(&x) {
		fmt.Printf("A:%d B:%s\n", x.A, x.B)
	}
	err = iter.Err()
	if err != nil {
		fmt.Printf("error iter: %s\n", err.Error())
		t.Fail()
	}
	iter.Close()
}
