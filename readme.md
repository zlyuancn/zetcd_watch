# 安全,可靠,稳定,使用简单的etcd键值监视模块

---

# 获得

`go get -u github.com/zlyuancn/zetcd_watch`

# 文档
[godoc](https://godoc.org/github.com/zlyuancn/zetcd_watch)

# 示例

```go
package main

import (
    "fmt"
    "go.etcd.io/etcd/clientv3"
    "github.com/zlyuancn/zetcd_watch"
)

func main() {
    c, _ := clientv3.NewFromURL("192.168.28.238:5460") // 连接
    manage := zetcd_watch.New(c)                       // 创建一个管理器
    w := manage.NewWatcher()                           // 创建一个监视器

    // 开始监视
    err := w.Start("", func(data *zetcd_watch.Data) {
        fmt.Printf("[%s]%s: %s\n", data.Type(), data.Key, data.Val())
    }, clientv3.WithPrefix())
    if err != nil {
        fmt.Println("错误", err)
    }
}
```

# 说明

## 此模块经过不停的断网、关闭etcd服务等严格测试, 能正常准确的监视key的每一次变动

## 解决官方etcd包的一个坑: 启用认证的etcd服务重启或断网重连后直接重新watch会token认证失败
