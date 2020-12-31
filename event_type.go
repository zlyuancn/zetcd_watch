/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/1/15
   Description :
-------------------------------------------------
*/

package zetcd_watch

// 事件类型
type EventType int

const (
	// 修改
	Change EventType = iota
	// 创建
	Create
	// 删除
	Delete
)

func (m EventType) String() string {
	switch m {
	case Create:
		return "Create"
	case Delete:
		return "Delete"
	}
	return "Change"
}
