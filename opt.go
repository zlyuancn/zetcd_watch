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
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// 默认重试等待时间
const DefaultRetryWaitTime = time.Second

// 选项
type options struct {
	wg            sync.WaitGroup
	ctx           context.Context  // 基础上下文, 用于通知结束
	client        *clientv3.Client // 官方的etcd客户端
	retryWaitTime time.Duration    // 重试等待时间
}

func newOptions(etcdClient *clientv3.Client) *options {
	return &options{
		ctx:           context.Background(),
		client:        etcdClient,
		retryWaitTime: DefaultRetryWaitTime,
	}
}

func (opt *options) apply(opts ...Option) {
	for _, o := range opts {
		o(opt)
	}
}

type Option func(opts *options)

// 设置上下文
func WithContext(ctx context.Context) Option {
	return func(opts *options) {
		if ctx == nil {
			ctx = context.Background()
		}
		opts.ctx = ctx
	}
}

// 设置监视断开重试等待时间
func WithRetryWaitTime(interval time.Duration) Option {
	return func(opts *options) {
		if interval <= 0 {
			interval = DefaultRetryWaitTime
		}
		opts.retryWaitTime = interval
	}
}
