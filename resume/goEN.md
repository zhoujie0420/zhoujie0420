# Jie Zhou

## Personal Information

* Phone: 18805268917 &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&ensp;&ensp; Email: jiezhou8917@qq.com
* Major: Computer Science &emsp;&emsp;&emsp;&emsp;&emsp;&nbsp; Position: Go Backend Engineer
* GitHub: https://github.com/zhoujie0420

***

## Skills

* Solid grasp of data structures and algorithms, with experience designing algorithms for complex real-world scenarios
* Proficient in Go, familiar with internals including map, slice, channel, GMP scheduling, GC, and memory escape analysis
* Experienced with Gin web framework; familiar with Argo Workflows for scheduled task orchestration; hands-on experience building agent workflows with the Eino framework
* Proficient in MySQL, MongoDB, and Redis with SQL tuning and index design experience
* Familiar with Kafka for async decoupling and traffic shaping scenarios
* Working knowledge of Linux, Docker, Kubernetes, and CI/CD pipelines
* Skilled in AI-native development tools (Claude Code, Cursor); experienced in writing high-quality prompts, maintaining project-level CLAUDE.md, and optimizing context configurations
* Solid understanding of cross-border payment business flows (clearing, settlement, acquiring, reconciliation, funding), with hands-on experience implementing fund consistency and regulatory compliance

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

### CardInfoLink — Cross-Border Payment Platform &emsp;&emsp; Role: Backend Engineer

A cross-border payment service provider operating across 10+ countries and regions with million-level daily transactions, covering acquiring, clearing, settlement, reconciliation, and funding modules. Primarily responsible for development and maintenance of clearing, settlement, and reconciliation modules.

* Participated in on-call rotation for the clearing & settlement module, diagnosing and mitigating production incidents to ensure fund accuracy and system stability for millions of daily transactions
* Resolved long-running card-organization BIN matching jobs by combining goroutine sharding with a greedy merge strategy, reducing job execution time from 10h to 30min
* Refactored the downstream service call chain with errgroup and DAG topological sorting to parallelize independent calls while preserving dependencies, reducing core API response time from 1.8s to 400ms
* Converted the reporter service into a resident process with a scheduling mechanism compatible with the existing scheduler, supporting multiple scheduling modes side-by-side and broadening file job flexibility
* Implemented end-to-end encryption for files to meet cross-border payment compliance requirements and maintained an internal decryption tool to lower manual processing cost
* Addressed slow end-of-day report jobs caused by full-table scans by introducing an intermediate aggregation layer with batched pre-aggregation, cutting report generation time by 70%
* Maintained project-level CLAUDE.md, encoding business rules for funding-file generation (amount precision, field specs, fund flow) as AI constraints, significantly reducing manual review cost and compliance risk in AI-generated code
* Built a reconciliation testing tool on top of the Eino framework that auto-generates edge-case scenarios and drives the full transaction lifecycle (initiation → clearing → reconciliation validation), covering one-sided entries, amount discrepancies, and cross-day delays, raising unit test coverage of core reconciliation modules to 99%

### Hang Seng Onboarding Platform &emsp;&emsp; Role: Backend Engineer (Core Developer)

Merchant onboarding platform for Hang Seng partners, covering the end-to-end flow from document submission, compliance review, agreement generation, rate configuration, to downstream payment-channel account provisioning — the front-end stage of the payment transaction chain.

* Designed and implemented a generic workflow engine abstracting node-level DAG transitions; submission, primary/secondary review, channel account creation, and activation steps are all configurable, allowing new nodes to be added without code changes
* Built an SLA timeout alerting mechanism matching SLA configs across multi-level dimensions (FunctionMenu, TaskType, etc.), computing ExpectedDate and Priority in real time from task creation timestamps to monitor onboarding timeliness
* Captured full-lifecycle audit trails via HistoryEntry, recording the operator, timestamp, and before/after values for every node entry, approval, rollback, and form change, meeting audit traceability and compliance requirements
* Abstracted hardcoded validation rules into JSON-configurable rules, supporting differentiated strategies per channel/institution and allowing new validation rules to be added without code changes

### Mogujie — E-commerce Main Site &emsp;&emsp; Role: Backend Engineer (Intern)

Main-site transaction-path development and shared scaffolding maintenance, covering order, refund, and rate-limiting modules with hundred-thousand-level daily orders.

* Delivered a subscription-purchase feature by reusing forward/reverse order logic, adding sub-order attachment and zero-cost purchase flows to reduce repetitive business wiring
* Implemented interface-level rate limiting using a leaky-bucket algorithm combined with AOP to cushion traffic bursts during peak events
* Applied distributed locks to enforce idempotency on low-level interfaces, reducing upstream complexity for duplicate-request and message-redelivery handling
* Unified scaffolding code standards (exception handling, in-house framework integration, etc.), cutting new-application bootstrap cost by 80%
* Rebuilt data-aggregation APIs with async orchestration to enable page pre-loading and improve first-screen experience
* Integrated third-party merchant systems with pluggable extension points for forward/reverse flows, turning new-channel onboarding from full-rewrite into plug-and-play configuration
