## 编译前端

```shell
go-bindata -o=asset/static/static.go -pkg=static -prefix=client/dist/static client/dist/static/...
```
```shell
go-bindata -o=asset/index/index.go -pkg=index -prefix=client/dist client/dist/favicon.ico client/dist/index.html
```