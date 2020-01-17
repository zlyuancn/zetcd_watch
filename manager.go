/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/11
   Description :
-------------------------------------------------
*/

package zetcd_watch

import (
    "context"
    "sync"
    "sync/atomic"
    "time"

    "github.com/zlyuancn/zlog2"
    "go.etcd.io/etcd/clientv3"
)

// 默认重试等待时间
const DefaultRetryWaitTime = time.Second

type Manager struct {
    // 最顶级的上下文, 用于通知关闭创建的watcher
    ctx context.Context
    // 上下文的关闭函数
    cancel context.CancelFunc
    // 用于等待所有watcher结束
    wg sync.WaitGroup

    // 官方的etcd客户端
    c *clientv3.Client
    // 是否运行中
    run int32
    // 重试等待时间
    retry_wait_time time.Duration
    // 日志
    log Loger
}

// 创建一个监视管理器
func New(etcd_client *clientv3.Client, opts ...Option) *Manager {
    ctx, cancel := context.WithCancel(context.Background())
    m := &Manager{
        ctx:    ctx,
        cancel: cancel,

        c:               etcd_client,
        run:             1,
        retry_wait_time: DefaultRetryWaitTime,
        log:             zlog2.DefaultLogger,
    }

    for _, o := range opts {
        o(m)
    }

    return m
}

// 是否运行中
func (m *Manager) IsRun() bool {
    return atomic.LoadInt32(&m.run) == 1
}

// 关闭并停止所有创建的Watcher
// 注意, 管理器不会主动关闭etcd客户端
func (m *Manager) Stop() {
    if atomic.CompareAndSwapInt32(&m.run, 1, 0) {
        m.cancel()
        m.wg.Wait()
    }
}

// 创建一个Watcher
func (m *Manager) NewWatcher() *Watcher {
    return newWatcher(m)
}

// 开启一个watcher, 所有的watcher在开始监视的时候必须调用它
func (m *Manager) startWatcher(w *Watcher) {
    m.wg.Add(1)
}

// 关闭一个watcher, 所有的watcher在停止监视的时候必须调用它
func (m *Manager) closeWatcher(w *Watcher) {
    m.wg.Done()
}
