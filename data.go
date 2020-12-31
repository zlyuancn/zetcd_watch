/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/11
   Description :
-------------------------------------------------
*/

package zetcd_watch

import (
	"github.com/zlyuancn/zstr"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type Kv struct {
	// 值
	Value zstr.String
	// 键
	Key string
	// 当前版本
	Version int64
	// 创建版本
	CreateRevision int64
	// 修订版本
	ModRevision int64
	// 租约id
	Lease int64
}

func makeKvFromRaw(kv *mvccpb.KeyValue) *Kv {
	return &Kv{
		Value:          zstr.String(kv.Value),
		Key:            string(kv.Key),
		Version:        kv.Version,
		CreateRevision: kv.CreateRevision,
		ModRevision:    kv.ModRevision,
		Lease:          kv.Lease,
	}
}

// 监视数据
type Data struct {
	Type  EventType
	*Kv       // 当前值
	OldKv *Kv // 原始值
}

func newData(e *clientv3.Event) *Data {
	m := &Data{
		Kv: makeKvFromRaw(e.Kv),
	}

	if e.PrevKv != nil {
		m.OldKv = makeKvFromRaw(e.PrevKv)
	}

	if e.Type == clientv3.EventTypeDelete {
		m.Type = Delete
	} else if m.CreateRevision == m.ModRevision {
		m.Type = Create
	} else {
		m.Type = Change
	}
	return m
}
