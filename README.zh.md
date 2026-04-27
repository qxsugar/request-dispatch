# Request Dispatch

一个强大的 Traefik 中间件插件，基于请求标记进行智能路由，支持灰度发布、受控功能发布和系统升级。

## 功能特性

- **智能请求路由** — 根据 HTTP 头将请求路由到不同的后端服务
- **灰度发布** — 逐步将新版本发布给用户子集
- **受控功能发布** — 在全量发布前，先向特定用户组测试新功能
- **负载均衡** — 为单个标记随机分配请求到多个后端主机
- **易于配置** — 基于 YAML 的简单配置，与 Traefik 无缝集成

## 安装

在 Traefik 配置中添加插件：

```yaml
experimental:
  plugins:
    request-dispatch:
      moduleName: github.com/qxsugar/request-dispatch
      version: v1.0.1
```

## 配置

### 基础设置

```yaml
http:
  routers:
    api:
      rule: host(`api.example.com`)
      service: production-service
      entryPoints:
        - web
      middlewares:
        - request-dispatch

  services:
    production-service:
      loadBalancer:
        servers:
          - url: "https://prod.api.example.com"

  middlewares:
    request-dispatch:
      plugin:
        request-dispatch:
          logLevel: INFO
          markHeader: X-Dispatch-Mark
          markHosts:
            canary:
              - https://canary-v2.api.example.com
              - https://canary-v2-backup.api.example.com
            staging:
              - https://staging.api.example.com
```

### 配置选项

| 选项           | 类型     | 说明                                    |
|--------------|--------|---------------------------------------|
| `logLevel`   | string | 日志级别：`DEBUG`、`INFO` 或 `ERROR`（不区分大小写） |
| `markHeader` | string | 用于检查路由标记的 HTTP 头名称                    |
| `markHosts`  | map    | 标记值到后端 URL 列表的映射                      |

## 使用

### 基础请求路由

发送带有标记头的请求以将其路由到特定后端：

```bash
# 路由到 canary 后端
curl -H "X-Dispatch-Mark: canary" https://api.example.com/users

# 路由到 staging 后端
curl -H "X-Dispatch-Mark: staging" https://api.example.com/users

# 路由到生产环境（默认）
curl https://api.example.com/users
```

### 负载均衡

为标记配置多个主机时，插件会随机选择其中一个：

```yaml
markHosts:
  canary:
    - https://canary-v2.api.example.com
    - https://canary-v2-backup.api.example.com  # 请求随机分配
```

### 灰度发布示例

```yaml
markHosts:
  v2:
    - https://api-v2.example.com
  v1:
    - https://api-v1.example.com
```

将 10% 的用户路由到 v2 进行测试：

```bash
# 10% 的请求
curl -H "X-Dispatch-Mark: v2" https://api.example.com/...

# 剩余 90% 的请求（默认路由）
curl https://api.example.com/...
```

## 开发

### 构建

```bash
go build
```

### 测试

```bash
make test          # 运行所有测试并显示覆盖率
make test -v       # 以详细模式运行测试
go test -v -run TestDispatch ./...  # 运行特定测试
```

### 代码检查

```bash
make lint
```

### 依赖管理

```bash
make vendor        # 下载依赖到 vendor 目录
make clean         # 删除 vendor 目录
```

## 架构

插件拦截 HTTP 请求并执行以下操作：

1. 检查配置的标记头
2. 如果头值匹配已配置的标记，随机选择一个后端主机
3. 将请求代理到选定的后端
4. 如果未找到标记或缺少头，回退到默认路由

### 线程安全

随机主机选择由互斥锁保护，确保在并发负载下的线程安全操作。

### 错误处理

如果反向代理失败，请求会自动回退到默认处理器，确保优雅降级。

## 日志

插件支持三个日志级别：

- **DEBUG** — 关于请求路由和代理操作的详细信息
- **INFO** — 插件操作的一般信息
- **ERROR** — 仅错误消息

配置中的日志级别不区分大小写。

## 许可证

详见 LICENSE 文件。
