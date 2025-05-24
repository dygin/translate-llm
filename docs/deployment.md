# AI 翻译系统部署文档

## 目录
- [系统要求](#系统要求)
- [部署架构](#部署架构)
- [前置准备](#前置准备)
- [部署步骤](#部署步骤)
- [监控系统](#监控系统)
- [日志系统](#日志系统)
- [CI/CD 配置](#cicd-配置)
- [常见问题](#常见问题)

## 系统要求

### 硬件要求
- CPU: 4核以上
- 内存: 8GB以上
- 存储: 50GB以上

### 软件要求
- Kubernetes 1.20+
- Docker 20.10+
- Helm 3.0+
- kubectl 1.20+

## 部署架构

系统采用微服务架构，主要包含以下组件：

1. **前端服务**
   - Vue 3 + TypeScript
   - Nginx 作为 Web 服务器
   - 静态资源 CDN 加速

2. **后端服务**
   - Go 1.21+
   - MySQL 8.0+
   - Redis 6.0+

3. **监控系统**
   - Prometheus
   - Grafana
   - Node Exporter

4. **日志系统**
   - Elasticsearch
   - Kibana
   - Filebeat

## 前置准备

1. **创建命名空间**
```bash
kubectl create namespace ai-translate
kubectl config set-context --current --namespace=ai-translate
```

2. **创建存储类**
```bash
kubectl apply -f k8s/storage-class.yaml
```

3. **创建密钥**
```bash
# 创建数据库密钥
kubectl create secret generic db-secrets \
  --from-literal=username=root \
  --from-literal=password=your-password

# 创建 Redis 密钥
kubectl create secret generic redis-secrets \
  --from-literal=password=your-redis-password
```

## 部署步骤

1. **部署数据库**
```bash
kubectl apply -f k8s/mysql-deployment.yaml
kubectl apply -f k8s/mysql-service.yaml
```

2. **部署 Redis**
```bash
kubectl apply -f k8s/redis-deployment.yaml
kubectl apply -f k8s/redis-service.yaml
```

3. **部署后端服务**
```bash
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/backend-service.yaml
```

4. **部署前端服务**
```bash
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/frontend-service.yaml
```

5. **部署 Ingress**
```bash
kubectl apply -f k8s/ingress.yaml
```

## 监控系统

### 部署监控组件
```bash
# 部署 Prometheus
kubectl apply -f k8s/monitoring/prometheus-config.yaml

# 部署 Grafana
kubectl apply -f k8s/monitoring/grafana-config.yaml
```

### 访问监控面板
```bash
# 端口转发
kubectl port-forward svc/grafana 3000:3000

# 访问 Grafana
# URL: http://localhost:3000
# 用户名: admin
# 密码: admin123
```

### 配置数据源
1. 登录 Grafana
2. 进入 Configuration > Data Sources
3. 添加 Prometheus 数据源
4. URL 设置为: http://prometheus:9090

### 导入仪表板
1. 进入 Dashboards > Import
2. 导入以下仪表板:
   - Kubernetes 集群概览
   - 应用性能监控
   - 资源使用统计

## 日志系统

### 部署日志组件
```bash
# 部署 ELK Stack
kubectl apply -f k8s/logging/elk-config.yaml
```

### 访问日志面板
```bash
# 端口转发
kubectl port-forward svc/kibana 5601:5601

# 访问 Kibana
# URL: http://localhost:5601
```

### 配置日志索引
1. 登录 Kibana
2. 进入 Stack Management > Index Patterns
3. 创建索引模式: filebeat-*
4. 选择时间字段: @timestamp

### 创建日志可视化
1. 进入 Visualize
2. 创建以下可视化:
   - 应用错误日志统计
   - 请求响应时间分布
   - 资源使用趋势

## CI/CD 配置

### GitHub Actions 配置

1. **添加 Secrets**
在 GitHub 仓库设置中添加以下 Secrets:
- `KUBE_CONFIG`: Kubernetes 集群的 kubeconfig 文件内容
- `DOCKER_USERNAME`: Docker Hub 用户名
- `DOCKER_PASSWORD`: Docker Hub 密码

2. **配置工作流**
- 推送代码到 main 分支将自动触发部署
- 创建 Pull Request 将触发测试
- 合并到 main 分支将触发构建和部署

### 手动触发部署
```bash
# 构建镜像
docker build -t your-registry/ai-translate-frontend:latest ./frontend
docker build -t your-registry/ai-translate-backend:latest ./backend

# 推送镜像
docker push your-registry/ai-translate-frontend:latest
docker push your-registry/ai-translate-backend:latest

# 更新部署
kubectl set image deployment/frontend frontend=your-registry/ai-translate-frontend:latest
kubectl set image deployment/backend backend=your-registry/ai-translate-backend:latest
```

## 常见问题

### 1. 服务无法访问
检查以下内容：
- 服务是否正常运行：`kubectl get pods`
- 服务端口是否正确：`kubectl get svc`
- Ingress 配置是否正确：`kubectl get ingress`

### 2. 数据库连接失败
检查以下内容：
- 数据库 Pod 状态：`kubectl get pods -l app=mysql`
- 数据库日志：`kubectl logs -l app=mysql`
- 数据库密钥是否正确：`kubectl get secret db-secrets`

### 3. 监控数据不显示
检查以下内容：
- Prometheus 状态：`kubectl get pods -l app=prometheus`
- Prometheus 日志：`kubectl logs -l app=prometheus`
- 服务监控注解是否正确

### 4. 日志收集异常
检查以下内容：
- Filebeat 状态：`kubectl get pods -l app=filebeat`
- Filebeat 日志：`kubectl logs -l app=filebeat`
- Elasticsearch 状态：`kubectl get pods -l app=elasticsearch`

### 5. CI/CD 部署失败
检查以下内容：
- GitHub Actions 日志
- 镜像构建日志
- Kubernetes 部署日志：`kubectl describe deployment frontend` 