/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/16
   Description :
-------------------------------------------------
*/

package zetcd_watch

import (
    "context"
    "errors"
    "fmt"
    "sync/atomic"
    "time"

    "go.etcd.io/etcd/clientv3"
)

// 监视键值改变的函数
type ObserverFunc func(data *Data)

var ManagerIsClose = errors.New("manager已关闭")
var WatcherIsRun = errors.New("这个watcher正在运行")

// 观察者
type Watcher struct {
    manager *Manager
    done    chan struct{}
    run     int32
}

func newWatcher(manager *Manager) *Watcher {
    w := &Watcher{
        manager: manager,
        done:    make(chan struct{}),
    }
    return w
}

// 开始监视
// 如果你需要监视一个key前缀, 请设置 clientv3.WithPrefix() 选项
func (m *Watcher) Start(key string, fn ObserverFunc, opts ...clientv3.OpOption) error {
    if fn == nil{
        panic("ObserverFunc is nil")
    }

    if !m.manager.IsRun() {
        return ManagerIsClose
    }

    if !atomic.CompareAndSwapInt32(&m.run, 0, 1) {
        return WatcherIsRun
    }

    m.manager.startWatcher(m)

    m.run = int32(1)
    ctx, cancel := context.WithCancel(m.manager.ctx)

    go func() {
        select {
        case <-ctx.Done():
            atomic.StoreInt32(&m.run, 0)
            <-m.done
        case <-m.done:
            atomic.StoreInt32(&m.run, 0)
            cancel()
        }
    }()

    var err error
    m.manager.log.Debug(fmt.Sprintf(`监视: "%s"`, key))
    for atomic.LoadInt32(&m.run) == 1 {
        // 强制要求切换token, 如果没有它, 在带认证的etcd服务重启后将会发生invalid token问题
        _, err = m.manager.c.Get(ctx, "/")
        if err != nil {
            cancel()
            break
        }

        err = m.watch(ctx, key, fn, opts...)
        if err != nil {
            m.manager.log.Warn(fmt.Errorf(`监视错误: "%s": %s`, key, err))
        }

        if atomic.LoadInt32(&m.run) == 0 {
            break
        }
        time.Sleep(m.manager.retry_wait_time)
        m.manager.log.Info(fmt.Sprintf(`重试监视: "%s"`, key))
    }

    m.manager.log.Debug(fmt.Sprintf(`停止监视: "%s"`, key))
    m.done <- struct{}{}

    m.manager.closeWatcher(m)
    return err
}

// 是否运行中
func (m *Watcher) IsRun() bool {
    return atomic.LoadInt32(&m.run) == 1
}

// 停止监视
func (m *Watcher) Stop() {
    if m.IsRun() {
        m.done <- struct{}{}
        <-m.done
    }
}

func (m *Watcher) watch(ctx context.Context, key string, fn ObserverFunc, opts ...clientv3.OpOption) error {
    ch := m.manager.c.Watch(ctx, key, opts...)
    for v := range ch {
        if v.Err() != nil {
            return v.Err()
        }
        for _, e := range v.Events {
            fn(newData(e))
        }
    }
    return nil
}
