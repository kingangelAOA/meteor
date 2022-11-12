### 通用化测试框架平台

测试平台的能力是可扩展的, 能力取决于框架中插件的能力, 此平台支持插件化开发, 然后编排加载插件的节点来提供测试能力, 假如插件支持前端自动化测试, 那么这个平台自然而然的支持前端自动化测试.以此类推

[![Watch the video](https://github.com/kingangelAOA/meteor/blob/main/doc/演示图片.png)]


#### 项目结构

```shell
/engine 编排node后,运行此编排的引擎
/plugin 插件目录

/web 平台代码
/web/client 平台前端代码
...

```

#### 插件
插件通过GRPC来进行通信, 所以理论上只要是支持GRPC的语言,都可以支持

支持: go, python
待支持: java, javascript

/plugin/python python 插件路径, 平台运行前, 运行 pip3 install -r req.txt

#### 编译平台前端
进入 meteor/web/client

```shell
1.npm install
2. npm run build
```

#### 编译平台后端
```shell
1. go-bindata -o=asset/static/static.go -pkg=static -prefix=client/dist/static client/dist/static/...
2. go-bindata -o=asset/index/index.go -pkg=index -prefix=client/dist client/dist/favicon.ico client/dist/index.html
3. go build web/main.go
```

#### 运行
```shell
./main
```
