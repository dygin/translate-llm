# AI翻译系统前端

基于 Vue 3 + TypeScript + Element Plus 开发的 AI 翻译系统前端项目。

## 功能特性

- 任务管理：创建、查看、更新、删除和重试任务
- 优先级管理：调整任务优先级，批量更新优先级
- 规则管理：创建和管理优先级规则、规则模板和规则组
- 响应式设计：适配不同屏幕尺寸
- 用户认证：登录和退出功能

## 技术栈

- Vue 3：渐进式 JavaScript 框架
- TypeScript：JavaScript 的超集，提供类型系统
- Element Plus：基于 Vue 3 的组件库
- Vue Router：Vue.js 的官方路由
- Axios：基于 Promise 的 HTTP 客户端
- Vite：下一代前端构建工具

## 开发环境要求

- Node.js >= 16.0.0
- npm >= 7.0.0

## 安装和运行

1. 安装依赖：

```bash
npm install
```

2. 启动开发服务器：

```bash
npm run dev
```

3. 构建生产版本：

```bash
npm run build
```

4. 预览生产构建：

```bash
npm run preview
```

## 项目结构

```
src/
├── api/          # API 请求
├── assets/       # 静态资源
├── components/   # 公共组件
├── router/       # 路由配置
├── styles/       # 全局样式
├── types/        # TypeScript 类型定义
├── utils/        # 工具函数
├── views/        # 页面组件
├── App.vue       # 根组件
└── main.ts       # 入口文件
```

## 代码规范

- 使用 ESLint 进行代码检查
- 使用 Prettier 进行代码格式化
- 遵循 Vue 3 组合式 API 风格指南

## 开发指南

1. 组件开发
   - 使用 `<script setup>` 语法
   - 使用 TypeScript 类型注解
   - 遵循组件命名规范

2. 样式开发
   - 使用 scoped CSS
   - 遵循 BEM 命名规范
   - 优先使用 Element Plus 的样式变量

3. API 开发
   - 使用 Axios 实例
   - 统一错误处理
   - 请求拦截器处理认证

## 部署

1. 构建项目：

```bash
npm run build
```

2. 将 `dist` 目录下的文件部署到 Web 服务器

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

[MIT](LICENSE) 