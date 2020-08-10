# **Go Entry Task**

## **1. [设计需求]**
- 实现一个用户管理系统，用户可以登陆，拉取和编辑他们的profiles.用户可以上传profile picture, 修改nickname。
- 分别实现处理HTTP请求的server和处理业务逻辑的TCP server.
- 用户账号信息必须存储在MySQL数据库
- 使用基于Auth/Session Token的鉴权机制
- 日志必须规范，需要分级，且单条可以追溯
- 业务逻辑必须通过grpc请求TCP服务器完成



## **2. [项目文件说明]**
项目主要源码都在src文件夹下。
- 第一个文件夹appendix下存储了各类小程序用于辅助系统测试和系统用户信息录入。以及服务器的功能测试
	1. 	crawl开头的程序，是用来对项目服务器的login接口进行http访问，从而获取指定用户数量的token,并写在对应的文件夹下生成csv文件，这些csv文件是用于后续使用jmeter测试接口性能必须的用户信息文件
	2. random_user.go是一个随机生成1000万用户账号和密码（此处为了方便，密码全设为123456）的程序，并写在了指定文件夹下。此处定义的函数生成200万用户，通过goroutine开了五个线程在同一文件夹生成了五个容量为200万的文件
	3. redis和TcpClient两个程序是单纯用于测试服务器和缓存连接的程序
	4. write2000Client和writeAllClients是用于记录指定数量的用户密码信息的程序，这是用于生成测试登陆接口的写文件的程序
	5. write_mysql是将此前用生成的1000万用户账号和密码信息写入mysql,由于一次性写太多会引起服务器崩溃，所以此处每次运行只写200万用户，需要运行五次，并且更改写入文件的信息

- config是配置文件，记录了mysql和redis的账号密码信息

- public是一个公用功能和服务接口实现的文件夹，config用于解析数据库配置文件，middleware中的logger是用于日志中间件的

- protos文件夹下存的是protobuf协议生成的服务接口和message对应的结构，用于支持grpc

- server文件夹下的程序是项目的主程序，分为两个服务器程序，http是指处理请求的http服务器，tcp是指处理业务逻辑的tcp服务器。
	1. http文件夹目录
		- api是用于处理http请求的，并且调用了grpc客户端服务
		- common文件夹下的go文件是用于连接tcp服务器的基础程序
		- connection文件夹下的go文件是用于构建基于grpc调用tcp服务器服务接口的连接器
		- datamodels是用于封装http请求参数和响应参数的结构体
		- middleware是用于鉴权的中间件文件夹
		- router是基于gin框架的路由映射器
		- web文件夹下的是前端静态资源（还在完善中）
    2. tcp文件夹目录
		- common文件夹下的go文件是用于各种基础数据处理和连接数据库的基本函数
		- datamodels下的文件是用于映射数据库实体的
		- repositories是数据连接仓库，用于对数据库进行sql操作，数据交互的结构体
		- service是服务，用于调用数据连接仓库进行数据获取并提供服务
	3. 运行项目时，必须两个文件夹下的main文件同时运行，开启服务器

## **3. [性能要求]**
- 200并发（固定用户），HTTP API QPS大于3000
- 200并发（随机用户），HTTP API QPS大于1000
- 2000并发（固定用户），HTTP API QPS大于1500
- 2000并发（随机用户），HTTP API QPS大于800



## **4. [开发工具]**
- IDE： goland
- 测试工具：jmeter、postman
- 开发语言：go1.14.6 darwin/amd64
- 数据库：mysql 8.0.21 for osx10.15 on x86_64
- 缓存：redis 6.0.6
- 框架：gin、grpc



## **6. [db设计]**
```
- create table user_tab(
	`user_id` bigint not null auto_increment,
	`username` varchar(255) not null,
	`password` varchar(255) not null,
	`nickname` varchar(255),
	`profile_picture` varchar(255),
	`created_time` timestamp not null default current_timestamp,
	`update_time` timestamp not null default current_timestamp on update current_timestamp,
	primary key (`user_id`)
	key `login_idx` (`username`,`password`)
) ENGINE=InnoDB auto_increment=1 DEFAULT CHARSET=utf8mb4
```

#### 1. 表格截图如下：
> ![截屏2020-08-09 下午2.51.55.png](https://i.loli.net/2020/08/09/NXzJcRWGiKoHOkM.png)

#### 2. 随机一个用户的信息截图：
> ![截屏2020-08-09 下午3.03.10.png](https://i.loli.net/2020/08/09/VPAz5LORpHTXfMb.png)

#### 3. 总用户数显示：
> ![截屏2020-08-09 下午3.05.08.png](https://i.loli.net/2020/08/09/ULGJfeN9A3Rdp7T.png)



## **7. [接口设计]**
1. 根据需求，需要有登陆，上传，修改三个功能，因此设计LogIn, ModifyNickname, UploadProfile三个主要功能
2. LogIn接口必须搭配鉴权中间件对用户进行认证，同时根据鉴权结果选择数据库和缓存的读写存取
3. ModifyNickname接口在接受参数前必须先进行鉴权判定用户处于登陆状态，然后才能进行修改
4. UploadProfile接口类似修改nickname的功能
5. 接口分为前端api用来接收http请求，以及后端service接口用来处理业务逻辑。两者通过grpc框架进行交互

### - Login
    以用户名和密码作为输入参数，进入后台服务进行查询，返回json如下：
    1. 当用户是第一次登陆时，需要生成token:
![截屏2020-08-09 下午4.26.27.png](https://i.loli.net/2020/08/09/yX8GTUNwIrREgsc.png)

![截屏2020-08-09 下午4.38.33.png](https://i.loli.net/2020/08/09/jgL9AQ2lDN7KthB.png)

	2. 当用户已经登陆过一次，根据token可以返回json结果:
![截屏2020-08-09 下午4.30.40.png](https://i.loli.net/2020/08/09/CloKyz9fGLjgwO4.png)

![截屏2020-08-09 下午4.34.55.png](https://i.loli.net/2020/08/09/IhsglEDnfi7BWQX.png)


### - ModifyNickname
	以校验码token和需要更改的nickname作为输入参数，返回json如下：
	1. 当用户登陆后，可以进行nickname修改(必须有token):
![截屏2020-08-09 下午4.42.54.png](https://i.loli.net/2020/08/09/kQyGVXC7cAPvTW1.png)

![截屏2020-08-09 下午4.43.59.png](https://i.loli.net/2020/08/09/SLHeujmRYXEIhs6.png)
	
![截屏2020-08-09 下午4.44.54.png](https://i.loli.net/2020/08/09/Vz7RhKj6lnToQDE.png)
	
	2. 若用户的token已经过期（redis中的token过期），登陆失效:
![截屏2020-08-09 下午4.47.44.png](https://i.loli.net/2020/08/09/cAGHt9jKJwFMYWP.png)
	
	3. 若是无token传参，是不允许的操作（用户未登陆）：
![截屏2020-08-09 下午4.54.13.png](https://i.loli.net/2020/08/09/mkzZIG3tn91M7A4.png)

### - UploadProfile
	以校验码token和需要上传的profile picture作为输入参数，返回json如下：
	1. 当用户登陆后，可以上传图片（url地址）(必须有token):
![截屏2020-08-09 下午4.50.17.png](https://i.loli.net/2020/08/09/xamfoDsyXMCqeAP.png)

![截屏2020-08-09 下午4.43.59.png](https://i.loli.net/2020/08/09/SLHeujmRYXEIhs6.png)

![截屏2020-08-09 下午4.51.27.png](https://i.loli.net/2020/08/09/1LUoN6ORhagZEHI.png)
	
	2. 若用户的token已经过期（redis中的token过期），登陆失效:
![截屏2020-08-09 下午4.47.44.png](https://i.loli.net/2020/08/09/cAGHt9jKJwFMYWP.png)

	3. 若是无token传参，是不允许的操作：
![截屏2020-08-09 下午4.54.13.png](https://i.loli.net/2020/08/09/mkzZIG3tn91M7A4.png)


## **8. [鉴权系统设计]**

- 这个系统的鉴权是基于jwt框架的，jwt是一个可以根据时间自动生成校验码的中间件框架。
- 首先用户第一次登陆必须由tcp向mysql申请获取信息，当用户名与密码校验成功，返回的用户模型不是空模型，则由返回的用户id和用户名（用户名唯一）以及当前时间生成下图定义的负荷，与先前定义的签名作为token加密生成的依据
![截屏2020-08-09 下午4.57.51.png](https://i.loli.net/2020/08/09/KplJD4zCIFwW5ZB.png)
	
- 以后每次当前用户只要访问系统的功能，就会自带这个token,直到token过期。此处token和用户信息一同存入缓存redis,作为高并发情况下有助于加速访问的中间件，这个操作可以加快高并发情况下的访问，增强系统的鲁棒性


## **9. [缓存层设计]**
缓存用redis, redis是一个单线程复用的key-value型nosql数据库。使用redis可以有效减缓mysql的压力，在高并发下有助于系统访问。此处运用redis存储键为用户登陆后生成的token, 值为用户的信息。

## **10. [性能测试前期准备]**
1. 对于login接口，首先用appendix下的writeAllClients和write2000clients程序进行用户信息记录。每一行代表一个用户名和密码，用逗号作为间隔。
2. 对于modifyNickname和uploadProfile接口，由于需要鉴权，必须获取token进行接口测试。此处，有两个方法解决这个问题：
	1. 使用appendix下的crawl程序对服务器进行访问，根据指定用户数量生成记录token的文件，然后用token文件作为jmeter的测量参数
	2. jmeter自带了请求前置和后置处理器。当随机用户量为10000000时，写入文件会导致机器卡死，且写入时间过久，项目的token有效期设置为1小时，可能等文件写入完毕，token已经过期了。因此，此处用了jmeter的后置处理器，第一个请求为login,然后用正则表达式获取第一个请求中得到的token,作为第二个请求中header的参数，就可以完成鉴权接口的测试。
		1. 后置请求器设置
![截屏2020-08-09 下午6.14.58.png](https://i.loli.net/2020/08/09/xUjOq9h7FL43vWt.png)
		2. 参数csv文件设置以及后置请求token设置
![截屏2020-08-09 下午6.15.10.png](https://i.loli.net/2020/08/09/vuCKwHQ2TanNJht.png)
![截屏2020-08-09 下午6.15.19.png](https://i.loli.net/2020/08/09/pg9YktANZzJSdU1.png)
        3. 后置请求头和参数设置（此处以上传接口为例）
![截屏2020-08-09 下午6.15.30.png](https://i.loli.net/2020/08/09/xSQZtKXYnAmFHja.png)

![截屏2020-08-09 下午6.15.41.png](https://i.loli.net/2020/08/09/zeoGEyKQ4hk83dB.png)

## **11. [性能测试报告]**
1. 根据要求，使用jmeter进行了性能测试，此处吞吐量可作为QPS指标。

- LogIn接口测试
	1. 200并发固定用户性能测试结果：
![200固定并发.png](https://i.loli.net/2020/08/09/g78pXO9im6xNjTh.png)
	
	2. 200并发随机用户性能测试结果：
![200随机并发.png](https://i.loli.net/2020/08/09/WEwYsQPKi8BqIGf.png)

	3. 2000并发固定用户性能测试结果：
![2000固定并发二测.png](https://i.loli.net/2020/08/09/NVRK43wXde7jWCT.png)
	
	4. 2000并发随机用户性能测试结果：
![2000并发随机二测.png](https://i.loli.net/2020/08/09/9zDQbqMyuNEO6LZ.png)

- ModifyNickname接口测试
	1. 200并发固定用户性能测试结果：
![200并发固定用户二测.png](https://i.loli.net/2020/08/09/8HJtEX9KcVAorFM.png)
	
	2. 200并发随机用户性能测试结果：
![200并发随机.png](https://i.loli.net/2020/08/09/NnaMm3uoG4H1zYh.png)

	3. 2000并发固定用户性能测试结果：
![2000并发固定用户修改二测.png](https://i.loli.net/2020/08/09/2FqUfSPMxWYnEjk.png)
	
	4. 2000并发随机用户性能测试结果：
![2000并发随机用户.png](https://i.loli.net/2020/08/09/MdOxHPlTiVWa69R.png)

- UploadProfile接口测试
	1. 200并发固定用户性能测试结果：
![200固定用户二测.png](https://i.loli.net/2020/08/09/Jh7IbCzVWvyFjTn.png)
	
	2. 200并发随机用户性能测试结果：
![200随机用户并发二测.png](https://i.loli.net/2020/08/09/UHakqWPQ2uEtep4.png)

	3. 2000并发固定用户性能测试结果：
![2000固定用户并发二测.png](https://i.loli.net/2020/08/09/7XEfxI9DCFvQnjK.png)
	
	4. 2000并发随机用户性能测试结果：
![2000随机用户二测.png](https://i.loli.net/2020/08/09/ER2wsFvyu7H3l6Y.png)



## **12. [项目心得]**
> 经过这次独立项目的开发，我学到了不少知识，包括了go语言的基本语法和应用，go的web框架和grpc远程调用的基础知识。同样也学习了go语言的鉴权机制。但是，我也同样发现了自己的一些不足，比如对缓存redis, mysql在实际业务场景下的应用缺乏经验以及对前端的不熟悉。今后还需继续学习，在生产业务中加强自我的技术和经验。
