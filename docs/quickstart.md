# AI 翻译系统快速入门指南

## 快速开始

### 1. 克隆代码
```bash
git clone https://github.com/your-username/ai-translate.git
cd ai-translate
```

### 2. 使用 Docker Compose 快速部署
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f
```

### 3. 访问服务
- 前端界面：http://localhost:80
- API 文档：http://localhost:8000/swagger
- Grafana 监控：http://localhost:3000
- Kibana 日志：http://localhost:5601

## 开发环境设置

### 1. 前端开发
```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 运行测试
npm run test

# 构建生产版本
npm run build
```

### 2. 后端开发
```bash
cd backend

# 安装依赖
go mod download

# 运行测试
go test ./...

# 启动开发服务器
go run cmd/main.go
```

## 常用命令

### Docker 命令
```bash
# 构建镜像
docker build -t ai-translate-frontend ./frontend
docker build -t ai-translate-backend ./backend

# 运行容器
docker run -d -p 80:80 ai-translate-frontend
docker run -d -p 8000:8000 ai-translate-backend

# 查看日志
docker logs -f container_id
```

### Kubernetes 命令
```bash
# 查看所有资源
kubectl get all

# 查看 Pod 日志
kubectl logs -f pod_name

# 进入 Pod
kubectl exec -it pod_name -- /bin/sh

# 查看服务状态
kubectl describe service service_name
```

### 监控命令
```bash
# 查看 Prometheus 目标
curl http://localhost:9090/api/v1/targets

# 查看 Grafana 数据源
curl http://localhost:3000/api/datasources

# 查看 Elasticsearch 索引
curl http://localhost:9200/_cat/indices
```

## 配置说明

### 环境变量
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ai_translate

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password

# 应用配置
APP_PORT=8000
APP_ENV=development
APP_DEBUG=true
```

### 配置文件
```yaml
# config.yaml
server:
  port: 8000
  mode: debug

database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  name: ai_translate

redis:
  host: localhost
  port: 6379
  password: your_password

logging:
  level: info
  filename: logs/app.log
  maxsize: 100
  maxage: 7
  maxbackups: 10
```

## 故障排除

### 1. 服务无法启动
```bash
# 检查端口占用
netstat -tulpn | grep LISTEN

# 检查日志
tail -f logs/app.log

# 检查配置
cat config.yaml
```

### 2. 数据库连接失败
```bash
# 检查数据库状态
mysql -u root -p -e "SHOW STATUS;"

# 检查数据库连接
mysql -u root -p -e "SHOW PROCESSLIST;"
```

### 3. Redis 连接失败
```bash
# 检查 Redis 状态
redis-cli ping

# 检查 Redis 连接
redis-cli info clients
```

## 开发建议

1. **代码规范**
   - 使用 ESLint 和 Prettier 保持代码风格一致
   - 遵循 Go 标准代码规范
   - 编写单元测试和集成测试

2. **Git 工作流**
   - 使用 feature 分支开发新功能
   - 提交前运行测试
   - 使用语义化版本号

3. **性能优化**
   - 使用连接池
   - 实现缓存策略
   - 优化数据库查询

4. **安全建议**
   - 使用环境变量存储敏感信息
   - 实现请求限流
   - 启用 HTTPS
   - 定期更新依赖

## 获取帮助

- 查看详细文档：`docs/` 目录
- 提交 Issue：GitHub Issues
- 联系支持：support@example.com 