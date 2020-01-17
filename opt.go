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
)

type Option func(m *Manager)

// 设置重试等待时间
func WithRetryWaitTime(t time.Duration) Option {
    return func(m *Manager) {
        if t <= 0 {
            t = DefaultRetryWaitTime
        }
        m.retry_wait_time = t
    }
}

// 自定义日志工具
func WithLogger(log Loger) Option {
    return func(m *Manager) {
        m.log = log
    }
}
