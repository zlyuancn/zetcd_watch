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
	// 完成等待
	done chan struct{}
	// 运行标志
	run int32
	// 记录版本
	recordVer map[string]int64

	// 日志工具
	log Loger
	// 官方客户端
	c *clientv3.Client
	// 重试等待时间
	retryWaitTime time.Duration
}

func newWatcher(manager *Manager) *Watcher {
	w := &Watcher{
		manager:       manager,
		done:          make(chan struct{}),
		log:           manager.log,
		c:             manager.c,
		retryWaitTime: manager.retryWaitTime,
	}
	return w
}

// 开始监视
// 如果你需要监视一个key前缀, 请设置 clientv3.WithPrefix() 选项
func (m *Watcher) Start(key string, fn ObserverFunc, opts ...clientv3.OpOption) error {
	if fn == nil {
		m.log.Warn(fmt.Sprintf("<%s>:ObserverFunc is nil", key))
		panic("ObserverFunc is nil")
	}

	if !m.manager.IsRun() {
		return ManagerIsClose
	}

	if !atomic.CompareAndSwapInt32(&m.run, 0, 1) {
		return WatcherIsRun
	}

	m.recordVer = make(map[string]int64)
	m.run = int32(1)
	ctx, cancel := m.manager.startWatcher(m)

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
	m.log.Debug(fmt.Sprintf(`监视: "%s"`, key))
	for atomic.LoadInt32(&m.run) == 1 {
		// 强制要求切换token, 如果没有它, 在带认证的etcd服务重启后将会发生invalid token问题
		_, err = m.c.Get(ctx, "/")
		if err != nil {
			cancel()
			break
		}

		err = m.watch(ctx, key, fn, opts...)
		if err != nil {
			m.log.Warn(fmt.Errorf(`监视错误: "%s": %s`, key, err))
		}

		if atomic.LoadInt32(&m.run) == 0 {
			break
		}
		time.Sleep(m.retryWaitTime)
		m.log.Info(fmt.Sprintf(`重试监视: "%s"`, key))
	}

	m.log.Debug(fmt.Sprintf(`停止监视: "%s"`, key))
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
	ch := m.c.Watch(ctx, key, opts...)
	for v := range ch {
		if v.Err() != nil {
			return v.Err()
		}
		for _, e := range v.Events {
			data := newData(e)
			if data.ModRevision == m.recordVer[data.Key] {
				continue // 过滤掉重复的
			}

			m.recordVer[data.Key] = data.ModRevision
			fn(data)
		}
	}
	return nil
}
