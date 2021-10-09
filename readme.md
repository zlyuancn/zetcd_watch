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
	client, _ := clientv3.NewFromURL("127.0.0.1:2379")
	_ = zetcd_watch.NewWatcher(client).Watch("/a", func(data *zetcd_watch.Data) {
		fmt.Printf("[%s] %s = %s \n", data.Type, data.Key, data.Value)
	})
}
```

# 可靠

## 此模块经过不停的断网、关闭etcd服务等严格测试, 能精确的监视key的每一次变动, 保证每次变动时精确触发一次回调
