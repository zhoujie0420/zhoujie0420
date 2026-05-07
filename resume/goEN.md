# Jie Zhou

## Personal Information

* Phone: 18805268917 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&ensp;&ensp; Email: jiezhou8917@qq.com
* Major: Computer Science &emsp;&emsp;&emsp;&emsp;&emsp;&nbsp; Position: Go Backend Engineer
* GitHub: https://github.com/zhoujie0420

***

## Skills

* Solid understanding of data structures and algorithms with strong problem-solving skills
* Proficient in Go, familiar with internals including map, slice, channel, GMP model, GC, and memory escape analysis
* Experienced with Gin web framework; familiar with Argo Workflows for scheduled task orchestration
* Proficient in MySQL, MongoDB, and Redis with SQL tuning experience
* Familiar with RabbitMQ messaging and related development practices
* Hands-on experience with Linux, Docker, Kubernetes, and CI/CD pipelines
* Skilled in AI-native development tools (Claude Code, Cursor); experienced in writing high-quality prompts, maintaining project-level CLAUDE.md, and optimizing context configurations

***

## Work Experience

**CardInfoLink** &emsp;&emsp;&emsp;&emsp; Backend Engineer &emsp;&emsp;&emsp;&emsp; 2024.01 – 2026.04

**Mogujie (Hangzhou Juangua Network)** &emsp;&emsp;&emsp;&emsp; Backend Engineer (Intern) &emsp;&emsp;&emsp;&emsp; 2023.08 – 2024.01

**AsiaInfo Technologies** &emsp;&emsp;&emsp;&emsp; Backend Engineer (Intern) &emsp;&emsp;&emsp;&emsp; 2023.04 – 2023.07

***

## Education

2020.09 – 2024.06 &emsp;&emsp; Zhejiang Shuren University &emsp;&emsp; Computer Science &emsp;&emsp; Bachelor's Degree

Achievements: CET-6 &emsp;&emsp; Academic Scholarship &emsp;&emsp; RoboCom Robotics Programming &emsp;&emsp; ACM

***

## Project Experience

### CardInfoLink — Cross-Border Payment Platform

A cross-border payment service provider covering clearing, acquiring, and platform business modules.

* Participated in on-call for the clearing & settlement module, diagnosing and resolving production issues to ensure system stability
* Designed a card BIN long/short matching algorithm, leveraging goroutines and greedy optimization to reduce task execution time from 10h to 30min
* Optimized reporting tasks with intermediate aggregation, reducing execution time by 70%
* Implemented parallel orchestration of downstream services using errgroup with DAG topological sorting for dependent call chains, reducing API response time from 1.8s to 400ms
* Refactored the file reporter service into a resident process with a scheduling mechanism supporting concurrent execution alongside existing schedulers
* Implemented end-to-end file encryption and maintained internal platform decryption tools to improve on-call analysis efficiency
* Used AI to simulate extreme reconciliation scenarios (one-sided transactions, amount discrepancies, cross-day delays), generating a unit test suite with 99% coverage
* Maintained CLAUDE.md to enforce AI code generation constraints for payment file logic, achieving 99% business compliance in AI-generated code

### Mogujie — E-commerce Main Site

Participated in main site development and infrastructure scaffolding maintenance.

* Developed subscription purchase feature by reusing existing forward/reverse order logic, adding sub-order attachment and zero-cost purchase flows to reduce repetitive business operations
* Optimized rate limiting with a leaky bucket algorithm and AOP-based interface-level throttling
* Implemented distributed locks for interface idempotency, reducing upstream business logic complexity
* Improved application scaffolding code standards (exception handling, internal framework usage), reducing new application setup cost by 80%
* Refactored data fetching interfaces using async orchestration for page preloading, improving user experience
* Integrated third-party applications with pluggable architecture reusing main site forward/reverse logic

### Godis — Redis Protocol Implementation in Go

A lightweight Redis protocol library implemented in Go.

* Designed and implemented core Redis data structures (strings, hashes, lists) supporting efficient storage and retrieval
* Optimized data operations using Go's concurrency primitives for high-throughput scenarios
* Implemented Redis protocol parsing and command execution with accurate response handling
* Improved caching strategies and memory management for better response times and resource utilization
* Completed persistence functionality including RDB and AOF file strategies
