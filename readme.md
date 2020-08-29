使用 [zlutils](https://github.com/lun-zhang/zlutils) 实现一个普通的web服务的完整例子

## 跑起来
1. 如果报错说什么sum not found，那就用go env -w GOPROXY="direct"把代理关掉
2. 这里已将从consul获取配置的代码注释了，方便直接启动，后续想用consul可以本地搭建个，很简单    
3. 需要本地有mysql server和redis server
4. 需要执行sql目录下的db.sql创建数据库，以及book.sql创建表

## 生成swagger接口文档
用[zlswagger](https://github.com/lun-zhang/zlswagger) 在当前目录执行
```shell script
./zlswagger \
-dir="./" \
-title="示例标题" \
-desc="示例描述"
```
即可在当前目录得到一个swagger.json文件