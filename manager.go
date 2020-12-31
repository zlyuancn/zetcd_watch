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
	"sync/atomic"

	"go.etcd.io/etcd/clientv3"
)

type Manager struct {
	// 最顶级的上下文, 用于通知关闭创建的watcher
	ctx context.Context
	// 上下文的关闭函数
	ctxCancel context.CancelFunc

	// 是否运行中
	run int32

	// 选项
	opts *options
}

// 创建一个监视管理器
func New(etcdClient *clientv3.Client, opts ...Option) *Manager {
	opt := newOptions(etcdClient)
	opt.apply(opts...)

	m := &Manager{
		run:  1,
		opts: newOptions(etcdClient),
	}
	m.ctx, m.ctxCancel = context.WithCancel(opt.ctx)
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
		m.ctxCancel()
		m.opts.wg.Wait()
	}
}

// 创建一个Watcher
func (m *Manager) NewWatcher() *Watcher {
	return newWatcher(m.opts)
}
