# 监控和日志配置指南

## 监控系统

### Prometheus 配置

#### 1. 基本配置
```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
```

#### 2. 告警规则
```yaml
# alert.rules
groups:
- name: example
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
    for: 10m
    labels:
      severity: critical
    annotations:
      summary: High error rate
      description: Error rate is {{ $value }} for the last 5 minutes
```

#### 3. 监控指标
- 系统指标：
  - CPU 使用率
  - 内存使用率
  - 磁盘 I/O
  - 网络流量
- 应用指标：
  - 请求延迟
  - 错误率
  - 并发连接数
  - 队列长度

### Grafana 配置

#### 1. 数据源配置
```yaml
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
```

#### 2. 仪表板配置
- 系统概览
  - CPU 使用率
  - 内存使用率
  - 磁盘使用率
  - 网络流量
- 应用监控
  - 请求延迟
  - 错误率
  - 并发用户数
  - 队列状态
- 资源使用
  - Pod 资源使用
  - 节点资源使用
  - 存储使用情况

#### 3. 告警配置
- 告警规则
  - 高 CPU 使用率
  - 高内存使用率
  - 高错误率
  - 服务不可用
- 通知渠道
  - 邮件
  - Slack
  - Webhook

## 日志系统

### Elasticsearch 配置

#### 1. 基本配置
```yaml
cluster.name: "ai-translate"
network.host: 0.0.0.0
http.port: 9200
discovery.type: single-node
```

#### 2. 索引配置
```yaml
# 索引模板
{
  "index_patterns": ["logs-*"],
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1
  },
  "mappings": {
    "properties": {
      "@timestamp": { "type": "date" },
      "level": { "type": "keyword" },
      "message": { "type": "text" }
    }
  }
}
```

#### 3. 日志保留策略
```yaml
# 索引生命周期策略
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_age": "7d",
            "max_size": "50gb"
          }
        }
      },
      "delete": {
        "min_age": "30d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

### Filebeat 配置

#### 1. 输入配置
```yaml
filebeat.inputs:
- type: container
  paths:
    - /var/lib/docker/containers/*/*.log
  processors:
    - add_kubernetes_metadata:
        host: ${NODE_NAME}
```

#### 2. 输出配置
```yaml
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  index: "filebeat-%{+yyyy.MM.dd}"
```

#### 3. 日志处理
- 日志格式
  - JSON 格式
  - 时间戳
  - 日志级别
  - 服务名称
  - 请求 ID
- 日志过滤
  - 排除健康检查日志
  - 排除调试日志
  - 保留错误日志

### Kibana 配置

#### 1. 索引模式
- 创建索引模式：`filebeat-*`
- 时间字段：`@timestamp`
- 字段映射：
  - 日志级别
  - 服务名称
  - 请求 ID
  - 错误信息

#### 2. 可视化配置
- 日志统计
  - 按级别统计
  - 按服务统计
  - 按时间统计
- 错误分析
  - 错误类型分布
  - 错误趋势
  - 错误详情

#### 3. 仪表板配置
- 系统日志
  - 日志数量趋势
  - 错误日志统计
  - 服务日志分布
- 应用日志
  - 请求日志分析
  - 错误日志分析
  - 性能日志分析

## 最佳实践

### 1. 监控最佳实践
- 设置合理的告警阈值
- 使用多级告警
- 定期检查告警规则
- 保持监控数据可追溯

### 2. 日志最佳实践
- 使用结构化日志
- 设置合理的日志级别
- 实现日志轮转
- 定期清理旧日志

### 3. 性能优化
- 使用日志缓冲
- 实现日志压缩
- 优化索引配置
- 合理设置资源限制

### 4. 安全建议
- 启用认证和授权
- 加密敏感信息
- 限制访问权限
- 定期审计日志

## 故障排除

### 1. Prometheus 问题
```bash
# 检查 Prometheus 状态
kubectl get pods -l app=prometheus

# 查看 Prometheus 日志
kubectl logs -l app=prometheus

# 检查目标状态
curl http://localhost:9090/api/v1/targets
```

### 2. Elasticsearch 问题
```bash
# 检查 Elasticsearch 状态
kubectl get pods -l app=elasticsearch

# 查看 Elasticsearch 日志
kubectl logs -l app=elasticsearch

# 检查集群健康状态
curl http://localhost:9200/_cluster/health
```

### 3. Filebeat 问题
```bash
# 检查 Filebeat 状态
kubectl get pods -l app=filebeat

# 查看 Filebeat 日志
kubectl logs -l app=filebeat

# 检查 Filebeat 配置
kubectl get configmap filebeat-config -o yaml
```

## 维护建议

### 1. 日常维护
- 检查监控系统状态
- 查看告警历史
- 分析日志趋势
- 优化资源配置

### 2. 定期维护
- 更新监控规则
- 优化日志配置
- 清理历史数据
- 检查系统性能

### 3. 应急处理
- 建立应急预案
- 准备回滚方案
- 保持系统备份
- 记录故障处理 