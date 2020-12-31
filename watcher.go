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
	"sync/atomic"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// 监视键值改变的函数
type ObserverFunc func(data *Data)

var WatcherIsRun = errors.New("这个watcher正在运行")

// 观察者
type Watcher struct {
	*options

	// 完成等待
	done chan struct{}
	// 运行标志
	run int32
	// 记录版本
	recordVer map[string]int64
}

func newWatcher(opts *options) *Watcher {
	w := &Watcher{
		options: opts,
		done:    make(chan struct{}),
	}
	return w
}

// 开始监视
// 如果你需要监视一个key前缀, 请设置 clientv3.WithPrefix() 选项
// 如果你想要知道key在改变之前的数据, 请设置 clientv3.WithPrevKV() 选项
func (w *Watcher) Watch(key string, fn ObserverFunc, opts ...clientv3.OpOption) error {
	if fn == nil {
		panic("ObserverFunc is nil")
	}

	if !atomic.CompareAndSwapInt32(&w.run, 0, 1) {
		return WatcherIsRun
	}

	w.recordVer = make(map[string]int64)
	w.wg.Add(1)

	ctx, cancel := context.WithCancel(w.ctx)
	go func() {
		select {
		case <-ctx.Done():
			atomic.StoreInt32(&w.run, 0)
			<-w.done
		case <-w.done:
			atomic.StoreInt32(&w.run, 0)
			cancel()
		}
	}()

	var err error
	for w.IsRun() {
		_ = w.watch(ctx, key, fn, opts...)
		if !w.IsRun() {
			break
		}
		time.Sleep(w.retryWaitTime)
	}

	w.done <- struct{}{}
	w.wg.Done()
	return err
}

// 是否运行中
func (w *Watcher) IsRun() bool {
	return atomic.LoadInt32(&w.run) == 1
}

// 停止监视
func (w *Watcher) Stop() {
	if w.IsRun() {
		w.done <- struct{}{}
		<-w.done
	}
}

func (w *Watcher) watch(ctx context.Context, key string, fn ObserverFunc, opts ...clientv3.OpOption) error {
	ch := w.client.Watch(ctx, key, opts...)
	for v := range ch {
		if v.Err() != nil {
			return v.Err()
		}
		for _, e := range v.Events {
			data := newData(e)
			if data.ModRevision == w.recordVer[data.Key] {
				continue // 过滤掉重复的
			}

			w.recordVer[data.Key] = data.ModRevision
			fn(data)
		}
	}
	return nil
}
