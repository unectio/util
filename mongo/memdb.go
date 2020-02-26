package mongo

import (
	"fmt"
	"sync"
	"errors"
	"reflect"
	"strings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	maxLimit int = 1024
)

type MemDB struct {
	lock	sync.RWMutex
	cols	map[string]*MemCol
}

func GetMemDB() Session {
	return &MemDB{
		cols: make(map[string]*MemCol),
	}
}

func (_ *MemDB)Close() {
}

func (mdb *MemDB)Copy() Session {
	return mdb
}

func (mdb *MemDB)Collection(loc *Location) Collection {
	mdb.lock.RLock()
	defer mdb.lock.RUnlock()

	key := loc.Db + "/" + loc.Col
	mcol, ok := mdb.cols[key]
	if !ok {
		mcol = &MemCol{}
		mdb.cols[key] = mcol
	}

	return mcol
}

type memObj struct {
	m	map[string]interface{}
	js	[]byte
}

func makeObj(o interface{}) (*memObj, error) {
	var err error

	ret := &memObj{}

	ret.js, err = bson.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("error marshal: %s", err.Error())
	}

	err = bson.Unmarshal(ret.js, &ret.m)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal: %s", err.Error())
	}

	return ret, nil
}

type MemCol struct {
	lock	sync.RWMutex
	objs	[]*memObj
}

func (mc *MemCol)Find(q bson.M) Query {
	return newQ(q, mc)
}

func (mc *MemCol)Insert(o interface{}) error {
	obj, err := makeObj(o)
	if err != nil {
		return err
	}

	mc.lock.RLock()
	defer mc.lock.RUnlock()

	mc.objs = append(mc.objs, obj)
	return nil
}

func (mc *MemCol)del(i int) {
	l := len(mc.objs) - 1
	mc.objs[i] = mc.objs[l]
	mc.objs = mc.objs[:l]
}

func (mc *MemCol)RemoveId(id bson.ObjectId) error {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mq := newQ(bson.M{"_id": id}, mc)
	i, _ := mq.find()
	if i != -1 {
		mc.del(i)
		return nil
	}

	return notFound()
}

func (mc *MemCol)RemoveAll(q bson.M) error {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	mq := newQ(q, mc)
	i := 0
	for i < len(mc.objs) {
		o := mc.objs[i]
		if !mq.match(o) {
			i++
		} else {
			mc.del(i)
		}
	}

	return nil
}

func (mc *MemCol)Update(q bson.M, u interface{}) error {
	mq := newQ(q, mc)
	i, o := mq.find()
	if i != -1 {
		no, err := updateObj(o, u)
		if err != nil {
			return err
		}

		mc.objs[i] = no
		return nil
	}

	return notFound()
}

func (mc *MemCol)EnsureIndex(_ *mgo.Index) error {
	return nil
}

func updateObj(o *memObj, u interface{}) (*memObj, error) {
	uo, err := makeObj(u)
	if err != nil {
		return nil, err
	}

	hasAct := false
	hasPlain := false

	for k, _ := range uo.m {
		if k[0] == '$' {
			if hasPlain {
				return nil, errors.New("Mixed action/plain update")
			}
			hasAct = true
		} else {
			if hasAct {
				return nil, errors.New("Mixed plain/action update")
			}
			hasPlain = true
		}
	}

	if hasPlain {
		return uo, nil
	}

	no := dupObj(o)

	for act, arg := range uo.m {
		switch act {
		case "$set":
			err = setObj(no, arg)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("Unsupported $-sction " + act)
		}
	}

	no.js, err = bson.Marshal(o.m)
	if err != nil {
		return nil, err
	}

	return no, nil
}

func setObj(obj *memObj, vals interface{}) error {
	switch vals := vals.(type) {
	case map[string]interface{}:
		return setObjFromMap(obj, vals)
	}

	fmt.Printf("$set tries to set %v\n", reflect.TypeOf(vals))
	return errors.New("unknwon $set type")
}

func setObjFromMap(o *memObj, vs map[string]interface{}) error {
	for k, v := range vs {
		err := setField(k, v, o)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyObj(from *memObj, to interface{}) error {
	return bson.Unmarshal(from.js, to)
}

func dupObj(o *memObj) *memObj {
	ret := &memObj{}

	err := bson.Unmarshal(o.js, &ret.m)
	if err != nil {
		panic("double unmarshal error")
	}

	ret.js = o.js
	return ret
}

type MemQ struct {
	q	bson.M
	mc	*MemCol
	lim	int
}

func newQ(q bson.M, mc *MemCol) *MemQ {
	return &MemQ{
		q:	q,
		mc:	mc,
		lim:	maxLimit,
	}
}

func (q *MemQ)Iter() Iter {
	return &MemIter{q, 0, nil}
}

func (q *MemQ)Limit(n int) Query {
	q.lim = n
	return q
}

func (q *MemQ)Sort(f string) Query {
	/* Sorry :( */
	return q
}

func (q *MemQ)Count() (int, error) {
	return -1, errors.New("Not implemented")
}

func (q *MemQ)One(out interface{}) error {
	q.mc.lock.RLock()
	defer q.mc.lock.RUnlock()

	i, o := q.find()
	if i == -1 {
		return notFound()
	}

	return copyObj(o, out)
}

func (q *MemQ)find() (int, *memObj) {
	for i, o := range q.mc.objs {
		if q.match(o) {
			return i, o
		}
	}
	return -1, nil
}

func equalS(a string, b interface{}) bool {
	bs, ok := b.(string)
	return ok && bs == a
}

func equalI(a int, b interface{}) bool {
	bi, ok := b.(int)
	return ok && bi == a
}

func equalOid(a bson.ObjectId, b interface{}) bool {
	bo, ok := b.(bson.ObjectId)
	return ok && bo == a
}

func equal(a interface{}, b interface{}) bool {
	switch a := a.(type) {
		case string:
			return equalS(a, b)
		case float64:
			return equalI(int(a), b)
		case int:
			return equalI(a, b)
		case bson.ObjectId:
			return equalOid(a, b)
	}

	fmt.Printf("unsupported type\n")
	return false
}

func setField(fname string, value interface{}, o *memObj) error {
	path := strings.Split(fname, ".")
	cur := o.m
	l := len(path)


	for i, n := range path {
		if i == l-1 {
			cur[n] = value
			return nil
		}

		x, ok := cur[n]
		if !ok {
			x = make(map[string]interface{})
			cur[n] = x
			continue
		}

		switch xt := x.(type) {
			case map[string]interface{}:
				cur = xt
			default:
				return errors.New("field conflict")
		}
	}

	return errors.New("field resolve fail")
}

func matchField(fname string, value interface{}, o *memObj) bool {
	path := strings.Split(fname, ".")
	cur := o.m
	l := len(path)


	for i, n := range path {
		x, ok := cur[n]
		if !ok {
			return false
		}
		if i == l-1 {
			return equal(x, value)
		}

		switch xt := x.(type) {
			case map[string]interface{}:
				cur = xt
			default:
				return false
		}
	}

	return false
}

func (q *MemQ)match(o *memObj) bool {
	for k, v := range q.q {
		if !matchField(k, v, o) {
			return false
		}
	}

	return true
}

type MemIter struct {
	q	*MemQ
	i	int
	er	error
}

func (i *MemIter)Err() error {
	return i.er
}

func (i *MemIter)Close() error {
	return i.er
}

func (i *MemIter)Next(out interface{}) bool {
	q := i.q
	if q.lim <= 0 {
		return false
	}

	for {
		if i.i == len(q.mc.objs) {
			return false
		}

		o := q.mc.objs[i.i]
		i.i++

		if q.match(o) {
			q.lim--
			err := copyObj(o, out)
			if err != nil {
				i.er = err
				return false
			}

			return true
		}
	}
}
