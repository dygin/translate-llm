# AI翻译字幕系统

基于GoFrame框架开发的AI翻译字幕系统，采用DDD架构设计和微服务架构。

## 项目结构

```
.
├── api                 # API接口定义
├── cmd                 # 命令行入口
├── internal            # 内部代码
│   ├── domain         # 领域模型
│   ├── application    # 应用服务
│   ├── infrastructure # 基础设施
│   └── interfaces     # 接口层
├── manifest           # 配置文件
├── resource           # 资源文件
└── web               # 前端代码
```

## 技术栈

### 后端
- GoFrame v2.6.1
- Redis
- MySQL
- Docker
- Kubernetes

### 前端
- Vue 3
- TypeScript
- Element Plus
- Vite

## 功能模块

1. 用户认证模块
2. 作品管理模块
3. 翻译管理模块
4. 调度管理模块
5. AI模型管理模块
6. 提示词管理模块
7. 权限管理模块

## 开发环境要求

- Go 1.21+
- Node.js 16+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 6.0+

## 快速开始

1. 克隆项目
2. 安装依赖
3. 配置环境变量
4. 运行项目

详细说明请参考各模块的文档。 