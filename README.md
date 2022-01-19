traefik 请求分发插件
------------------

根据请求mark做请求转发，实现灰度功能

## 配置

1. 启用插件配置

```yaml
experimental:
  plugins:
    requst-dispatch:
      moduleName: github.com/qxsugar/request-dispatch
      version: v1.0.1
```

2. 路由配置

```yaml
http:
  routers:
    api:
      rule: host(`api.test.cn`)
      service: svc
      entryPoints:
        - web
      middlewares:
        - requst-dispatch
  services:
    svc:
      loadBalancer:
        servers:
          - url: "http://prod.api.cn"
  middlewares:
    requst-dispatch:
      plugin:
        requst-dispatch:
          logLevel: DEBUG
          markHeader: X-Tag
          markHosts:
            alpha:
              - http://alpha.api.cn
              - http://alpha1.api.cn
            beta:
              - http://beta.api.cn
```

## 使用

如果请求头带有markHeader的参数，请求将会被分发到对应的地址

> http api.test.cn X-Tag:alpha -v # 请求会被分发到`http://alpha.api.cn`或`http://alpha1.api.cn`
>
> http api.test.cn X-Tag:beta -v # 请求会被转分发到`http://beta.api.cn`
>
> http api.test.cn X-Tag:whoami -v # 请求不会被分发