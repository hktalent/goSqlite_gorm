package common

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
)

// https://colobu.com/2017/10/11/badger-a-performant-k-v-store/
// https://juejin.cn/post/6844903814571491335
type KvDbOp struct {
	DbConn *badger.DB
}

var cache *KvDbOp

func NewKvDbOp() *KvDbOp {
	if nil != cache {
		return cache
	}
	r := KvDbOp{}
	r.Init("db/DbCache")
	cache = &r
	return cache
}

func (r *KvDbOp) Init(szDb string) error {
	opts := badger.DefaultOptions(szDb)
	db, err := badger.Open(opts)
	if nil != err {
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
func PutAny[T any](r *KvDbOp, key string, data T) {
	d, err := json.Marshal(data)
	if nil == err {
		r.Put(key, d)
	}
}

func GetAny[T any](r *KvDbOp, key string) T {
	var t1 T
	data, err := r.Get(key)
	if nil == err {
		json.Unmarshal(data, &t1)
		return t1
	}
	return t1
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
