/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
/////////////////////////////////////////////////////////////////////////////////

package mongo

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

type Foo struct {
	Id bson.ObjectId `bson:"_id"`
	A  int           `bson:"a"`
	B  string        `bson:"b"`
}

func TestMemdb(t *testing.T) {
	mdb := GetMemDB()

	col := mdb.Collection(&Location{"foo", "bar"})
	err := col.Insert(&Foo{Id: bson.NewObjectId(), A: 1, B: "one"})
	if err != nil {
		fmt.Printf("error insert: %s\n", err.Error())
		t.FailNow()
	}

	if err := col.Insert(&Foo{Id: bson.NewObjectId(), A: 2, B: "two"}); err != nil {
		t.Fatalf("error insert: %s\n", err)
	}

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
