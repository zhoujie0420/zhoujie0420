# 周杰

## 个人信息

* 手 机：18805268917 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&ensp;&ensp; 邮 箱：jiezhou8917@qq.com
* 专 业：计算机科学与技术 &emsp;&emsp;&emsp;&emsp;&emsp;&nbsp; 岗 位：Go开发工程师
* GitHub：https://github.com/zhoujie0420

***

## 专业技能

* 熟悉数据结构和算法，具备良好的算法思想和思维逻辑，熟悉高级数据结构的使用
* 熟悉Go语言，了解map、slice、channel、GMP模型、gc垃圾回收、内存逃逸等底层原理
* 熟悉 Gin Web 框架，了解 Argo Workflows 并用于定时任务调度与编排
* 熟悉常见数据库如MySQL、MongoDB、Redis，具备SQL调优经验
* 熟悉消息队列 RabbitMq 及其基本使用，拥有相关的开发经验。
* 了解并能使用Linux操作系统及Docker、K8S，熟悉CI/CD流程。
* 熟练运用 Claude Code、Cursor 等 AI 原生开发工具，擅长编写高质量 Prompt、维护项目级 CLAUDE.md 及优化上下文配置

***

## 工作经历

公 司 ：讯联数据&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;职位 ： 后端
&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;时间 ：2024.1-2026.4

公 司 ：蘑菇街（杭州卷瓜网络）&emsp;&emsp;职位 ： 后端（实习） &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;时间 ：2023.8-2024.1

公 司 ：亚信科技（中国） &emsp;&emsp;&emsp;&emsp;&emsp;职位 ：后端（实习） &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;时间
：2023.4-2023.7

***
## 教育经历
2020.09  -  2024.06 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;浙江树人大学   &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;计算机科学与技术   &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;本科   
在校经历： CET-6 &emsp;&emsp;&emsp;&emsp;学年奖学金  &emsp;&emsp;&emsp;&emsp;RoboCom 机器人开发编程   &emsp;&emsp;&emsp;ACM

***

## 项目经历

项 目 ：CardInfoLink
&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp; 

项目介绍：跨境支付服务商，涉及清算，收单，平台业务模块。

* 参与清结算模块的oncall工作，定位并解决线上问题，保障系统稳定运行
* 设计cardBin长短位匹配算法，结合协程、贪心算法优化合并效率，优化任务执行时间10h->30min
* 优化报表任务，针对报表数据进行中间统计，任务执行时间缩短70%
* 基于 errgroup 实现多下游服务并行编排，采用 DAG 拓扑排序思想处理有依赖的调用链，接口 RT 从 1.8s 降至 400ms
* 文件reporter服务优化为常驻进程，设计调度机制，支持和现有调度同时执行，拓展文件调度方式
* 文件数据的全流程加密，维护内部平台的文件解密功能，提高内部oncall分析效率
* 利用 AI 模拟极端对账场景（如单边账、金额差错、跨日延迟），生成了 99% 覆盖单元测试集
* 维护 CLAUDE.md,强制规范 AI 在生成请款文件逻辑时的规范约束,确保 AI 生成代码的业务合规性达 99%。

项 目 ：蘑菇街
&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;
 

项目描述：参与蘑菇街主站开发及基础脚手架维护。

* 周期购需求，复用原有的正逆向逻辑，新增子订单依附及零元购通用逻辑，简化业务重复操作
* 优化缓存限流策略，使用漏斗算法及Aop思想实现接口级限流
* 使用分布式锁实现底层接口的幂等性，减少上游业务的复杂逻辑的考虑。
* 优化应用脚手架的代码规范，包括异常处理，自研框架应用等，降低新建应用成本80%
* 使用异步编排的方式重构获取数据接口，实现页面的预加载，提高用户体验
* 对接第三方应用，复用主站正向逆向逻辑，实现新增三方的可插拔式开发


项 目 ：Godis - Go语言实现的Redis协议库
&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;

项目描述：Godis项目的开发和优化，实现一个轻量级的Redis协议库。
* 基于Go语言设计并实现Redis的核心数据结构，如字符串、哈希、列表等，支持高效的数据存储与读取。
* 使用Go语言的并发机制优化数据操作，提升系统在高并发下的处理能力。
* 实现了Redis协议的解析和处理逻辑，确保了命令的准确执行与响应。
* 通过改进缓存策略和内存管理，提升了系统的响应速度和资源利用率。
* 完成了系统的持久化功能，包括RDB和AOF文件持久化策略的实现。