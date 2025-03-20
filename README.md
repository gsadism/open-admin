# open-admin



>**约定**
>
>- [arg1 | arg2] ：可选参数

## 启动方式

~~~bash
go run main.go [server] # 使用系统默认配置方式启动

go run main.go [server] --conf=./application.yml --log=debug # 指定启动方式
~~~

- `--conf`: 制定系统配置文件
- `--log`: 控制台日志输出级别 [debug、info、warn、error]



## 环境变量

|          变量名          |                             描述                             |
| :----------------------: | :----------------------------------------------------------: |
|   **OPEN_ADMIN_ROOT**    |           项目运行根目录,系统自动设置,用户无需配置           |
|    **OPEN_ADMIN_ID**     |  服务唯一标识,常用于分布式环境。当用户不指定时,则会随机生成  |
| **OPEN_ADMIN_LOG_LEVEL** | 控制台日志输出级别,优先级：命令行参数>环境变量>默认配置.[debug, info,warn,error] |
|                          |                                                              |
|                          |                                                              |

