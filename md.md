# 周杰
## 个人信息 

* 性 别：男&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&ensp;年 龄：22  
* 手 机：18805268917 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&ensp;&ensp; 邮 箱：1747364257@qq.com    
* 专 业：计算机科学与技术 &emsp;&emsp;&emsp;&emsp;&emsp; 岗 位：Java开发工程师
* Github：https://github.com/zhoujie0420
***
## 专业技能

* 熟练掌握 java 相关知识 ，具备良好的面向对象的编程思想。
* 熟悉常见数据结构和算法 ，具备良好的算法思想和思维逻辑。
* 熟悉SpringMVC，Mybatis-Plus，SpringBoot等主流框架，能够使用框架完成restful接口编写。
* 熟悉编写SQL语句，具备一定的数据库表设计能力，熟悉索引，具备 sql 优化经验。
* 熟悉 Redis 数据类型的使用，实践缓存穿透，缓存击穿和缓存雪崩情况解决，了解其持久化策略。
* 了解常用开发工具的使用 如：Maven，Git ，Jmeter，Postman ，熟练使用开发工具 Idea。
* 了解并能使用Linux操作系统以及docker容器，有实际部署项目的经验.
* 了解消息队列 RabbitMQ 及其基本使用。
* 掌握 Vue 和开源 Ui 框架的前端开发。


***
## 教育经历

2020.09  -  2024.06 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;浙江树人学院   &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;计算机科学与技术   &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;本科   
在校经历： CET-4 &emsp;&emsp;&emsp;&emsp;&emsp;2021学年奖学金  &emsp;&emsp;&emsp;&emsp;&emsp;RoboCom 机器人开发者大赛编程设计赛省三   &emsp;&emsp;&emsp;&emsp;&emsp;蓝桥杯省三

***
## 实习经历

 公 司 ：杭州匡汇科技公司  &emsp;&emsp;&emsp;&emsp;&emsp;职位 ：Java 后端 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;时间 ：2023.1-2023.5

* 高效地完成各项任务，并与前端团队密切合作，进行联调工作，以确保项目顺利推进，提高工作效率和产出。
* 负责将新老项目的数据进行转移和整合。加快新版本项目的上线进度，编写大量视图，用于快速迁移和操作数据，确保项目顺利进行。
* 用户和角色模块的查询功能进行了封装处理，为团队成员提供了高效开发的工具，提高了整体开发效率。
* 针对异步下载接口，我进行了重构工作，并提供了查询条件的输入功能，提高接口的复用性和灵活性。

***
## 项目经历

项  目 ：Acmer Of Bots  &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;     项目角色：全栈

项目描述：基于SpringBoot+Vue3的本校ACM社团招新宣传的项目，用户利用自己的算法知识完成rating提升
* 前端采用Vue框架搭建，使用bootstrap组件库
* 使用Spring Security+jwt 实现用户登录，权限管理
* 采用微服务架构，完成匹配系统和Ai代码执行业务，高效实现系统对局完整性，提高游戏体验
* 使用WebSocket 长连接实现回合制游戏流程
* 使用 @scheduled 定时任务，分析对局情况，用户上线情况
* 运用 RabbitMq 完成存储对局记录，提高系统的性能和并发处理能力
* 使用docker完成项目的部署，利用云服务器完成项目的上线


项  目  ：店铺点评APP   &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;     项目角色：后端

项目描述：基于SpringBoot+Redis的店铺点评APP,实现了找店铺、写点评、看热评、点赞、关注Feed流的完整业务流程
* 使用Redis实现分布式Session解决集群间登录态同步的问题，利用双重拦截模式实现了登录的权限控制和状态刷新
* 基于静态ThreadLocal封装了线程隔离的全局上下文对象，便于请求内部存取用户信息
* 为提高ID的安全性和高可用性，利用时间戳和Redis的INCR命令构造了全局唯一ID生成器，支持每秒最多能生成
2^32个不同的ID
* 使用Redis对高频访问的信息进行缓存，降低数据库压力提高了80%的查询性能，解决了缓存穿透,雪崩、击穿问题。
* 使用 RabbitMq 实现数据库的批量操作，缓解数据库压力，提高性能 。
* 优惠券秒杀:使用Redis+Lua脚本实现库存预检，并通过Stream消息队列实现订单的异步创建，解决超卖问题实现一
人一单，性能提高50%
* 使用Redis ZSet数据结构存储用户的点赞信息并实现TopN点赞排行榜，使用Redis Set数据结构实现用户关注、共
同关注，使用Redis Geo实现地理位置的存储。
* 基于推模式实现关注Feed流，保证博主的更新内容及时推送给粉丝，减少用户访问等待的时间

