package common

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"log"
)

var cache1 = NewKvDbOp()

// https://colobu.com/2017/10/11/badger-a-performant-k-v-store/
// https://juejin.cn/post/6844903814571491335
type KvDbOp struct {
	DbConn *badger.DB
}

func NewKvDbOp() *KvDbOp {
	r := KvDbOp{}
	r.Init("db/DbCache")
	return &r
}

func (r *KvDbOp) Init(szDb string) error {
	opts := badger.DefaultOptions(szDb)
	db, err := badger.Open(opts)
	if nil != err {
		log.Println("Init k-v db 不能多个进程同时开启", err)
		return err
	}
	r.DbConn = db
	return nil
}

func (r *KvDbOp) Delete(key string) error {
	err := r.DbConn.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	return err
}

func (r *KvDbOp) Close() {
	r.DbConn.Close()
}

func (r *KvDbOp) Get(key string) (szRst []byte, err error) {
	err = r.DbConn.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		// val, err := item.Value()
		err = item.Value(func(val []byte) error {
			szRst = val
			return nil
		})
		return err
	})
	return szRst, err
}
func PutAny[T any](key string, data T) {
	d, err := json.Marshal(data)
	if nil == err {
		cache1.Put(key, d)
	}
}

func GetAny[T any](key string) (T, error) {
	var t1 T
	data, err := cache1.Get(key)
	if nil == err {
		json.Unmarshal(data, &t1)
		return t1, nil
	}
	return t1, err
}

func (r *KvDbOp) Put(key string, data []byte) {
	err := r.DbConn.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), data)
		if err == badger.ErrTxnTooBig {
			_ = txn.Commit()
		}
		return err
	})
	if err != nil {
	}
}
