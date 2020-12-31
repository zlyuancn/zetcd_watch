/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/16
   Description :
-------------------------------------------------
*/

package zetcd_watch

import (
	"time"

	"github.com/zlyuancn/zlog"
)

type Option func(m *Manager)

// 设置监视断开重试等待时间
func WithRetryWaitTime(t time.Duration) Option {
	return func(m *Manager) {
		if t <= 0 {
			t = DefaultRetryWaitTime
		}
		m.retryWaitTime = t
	}
}

// 自定义日志工具
func WithLogger(log Loger) Option {
	return func(m *Manager) {
		if log == nil {
			log = zlog.DefaultLogger
		}
		m.log = log
	}
}
