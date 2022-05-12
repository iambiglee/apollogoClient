# Apollo Go Client

本项目计划是做一个[Apollo](https://github.com/apolloconfig/apollo) 的GO的客户端。

如果计划在正式环境中使用Apollo， 并且有GO的客户端需要使用Apollo 作为配置中心，建议使用[agollo](https://github.com/apolloconfig/agollo)。

我会尽可能的将设计思路写出来。

1. 初始化配置文件
   agollo 启动的时候，先要去寻找配置了apollo.配置的内容信息如下

   ```json
   {
       "appId": "agollo-test",
       "cluster": "dev",
       "namespaceName": "testjson.json,testyml.yml",
       "ip": "http://106.54.227.205:8080",
       "releaseKey": ""
       "secret":"7c2ddeb1cd344b8b8db185b3d8641e7f"
   }
   ```

2. 连接配置文件