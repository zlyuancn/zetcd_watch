/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/16
   Description :
-------------------------------------------------
*/

package zetcd_watch

// 日志接口
type Loger interface {
    Debug(v ...interface{})
    Info(v ...interface{})
    Warn(v ...interface{})
}
