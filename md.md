# 周杰
## 个人信息 

 
* 手 机：18805268917 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&ensp;&ensp; 邮 箱：1747364257@qq.com    
* 专 业：计算机科学与技术 &emsp;&emsp;&emsp;&emsp;&emsp; 岗 位：Java开发工程师
* GitHub：https://github.com/zhoujie0420
***
## 专业技能

* 熟悉数据结构和算法 ，具备良好的算法思想和思维逻辑
* 了解 JUC，JVM 具备 Java 开发经验，熟悉 SpringBoot，SpringCloud，Mybatis-Plus等框架
* 熟悉 Mysql，具备一定的数据库表设计能力，熟悉索引，具备 sql 优化经验
* 熟悉 Redis 数据类型的使用，实践缓存穿透，缓存击穿和缓存雪崩情况解决，了解其持久化策略
* 了解消息队列RabbitMQ及其基本使用，拥有相关的开发经验。
* 了解常用开发工具的使用 如：Maven，Git ，Jmeter，Postman ，熟练使用开发工具 Idea
* 了解并能使用Linux操作系统以及docker容器，有实际部署项目的经验



***
## 教育经历

2020.09  -  2024.06 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;浙江树人学院   &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;计算机科学与技术   &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;本科   
在校经历： CET-4 &emsp;&emsp;&emsp;&emsp;2021学年奖学金  &emsp;&emsp;&emsp;&emsp;RoboCom 机器人开发编程设计   &emsp;&emsp;&emsp;蓝桥杯省三

***
## 实习经历

 公 司 ：亚信科技（中国）   &emsp;&emsp;&emsp;&emsp;&emsp;职位 ：Java 后端 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;时间 ：2023.5-2023.7

项目介绍：与国内三大运营商合作的ToB产品，提供不同平台的广告投送。
* 负责生产环境的订单库存的数据稳定，进行需求回滚。编写接口文档
* 根据产品需求，完善部分消息监听接口，实现异步接口的可靠，提高接口响应速度
* 负责用户视频观看功能模块，使用缓存加定时任务实现数据持久化，降低DB压力
  

 公 司 ：杭州匡汇科技公司  &emsp;&emsp;&emsp;&emsp;&emsp;职位 ：Java 后端 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;时间 ：2023.1-2023.5
 
 地址：http://www.grs.zju.edu.cn/
 
项目介绍：属于 OA 管理项目，为各高校研发教务系统，提供多角色应用功能，保障校园生活智能化

* 负责将新老项目的数据进行转移和整合。编写视图，迁移和操作数据，确保项目顺利进行
* 组装部分角色查询条件，编写开发文档，提高了团队开发效率
* 针对异步下载接口进行重构，提供查询条件的泛型接口，提高接口的复用性和灵活性

***
## 项目经历
项  目 ：Acmer Of Bots  &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;     项目角色：全栈

地址：https://zzac.online/

项目描述：基于SpringBoot+Vue3的本校ACM社团招新宣传的项目，用户利用自己的算法知识完成rating提升
* 使用Spring Security+jwt 实现用户登录，权限管理
* 隔离匹配系统和Ai代码执行服务，高效实现系统对局完整性，提高游戏体验
* 使用 WebSocket 长连接实现游戏同步,异步存储对局记录，提高系统的性能
* 使用 Quartz 异步任务实现对局分析及用户上线情况统计
* 使用 docker 完成项目的部署，利用云服务器完成项目的上线

项  目  ：店铺点评APP   &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;     项目角色：后端

项目描述：基于Redis的店铺点评APP,对接口进行重构，提高性能
* 使用Redis实现分布式Session解决集群间登录态同步的问题，利用双重拦截模式实现了登录的权限控制和状态刷新
* 基于静态ThreadLocal封装了线程隔离的全局上下文对象，便于请求内部存取用户信息
* 为提高ID的安全性和高可用性，利用时间戳和Redis的INCR命令构造了全局唯一ID生成器，支持每秒最多能生成
2^32个不同的ID
* 使用Redis对高频访问的信息进行缓存，降低DB压力提高了80%的查询性能，解决了缓存穿透,雪崩、击穿问题
* 使用 RabbitMq 实现数据库的批量操作，缓解DB压力，提高性能
* 优惠券秒杀:使用Redis+Lua脚本实现库存预检，并通过Stream消息队列实现订单的异步创建，解决超卖问题实现一
人一单，性能提高50%