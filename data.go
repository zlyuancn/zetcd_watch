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
)

// 监视数据
type Data struct {
    t EventType
    // 值
    zstr.String
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

func newData(e *clientv3.Event) *Data {
    kv := e.Kv
    m := &Data{
        String:         zstr.String(kv.Value),
        Key:            string(kv.Key),
        Version:        kv.Version,
        CreateRevision: kv.CreateRevision,
        ModRevision:    kv.ModRevision,
        Lease:          kv.Lease,
    }

    if e.Type == clientv3.EventTypeDelete {
        m.t = Delete
    } else if m.CreateRevision == m.ModRevision {
        m.t = Create
    } else {
        m.t = Change
    }
    return m
}

// 是否为修改
func (m *Data) IsChange() bool {
    return m.t == Change
}

// 是否为创建
func (m *Data) IsCreate() bool {
    return m.t == Create
}

// 是否为删除
func (m *Data) IsDelete() bool {
    return m.t == Delete
}

// 事件类型
func (m *Data) Type() EventType {
    return m.t
}
