# 安全,可靠,稳定,使用简单的etcd键值监视模块

---

# 获得

`go get -u github.com/zlyuancn/zetcd_watch`

# 文档
[godoc](https://godoc.org/github.com/zlyuancn/zetcd_watch)

# 示例

```
// 使用官方的etcd客户端库
c, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"127.0.0.1:2379"}, // 集群或单点地址
    DialTimeout: 5 * time.Second,
})
if err != nil {
    panic(err)
}
defer c.Close()

// 创建一个管理器
manage := zetcd_watch.New(c)

// 创建一个监视器
w := manage.NewWatcher()

// 开始监视
err = w.Start("", func(data *zetcd_watch.Data) {
    fmt.Printf("[%s]%s: %s\n", data.Type(), data.Key, data.Val())
}, clientv3.WithPrefix())
if err != nil {
    fmt.Println("错误", err)
}
```
