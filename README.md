traefik request dispatch plug-in
------------------

According to the request mark to do request dispatch, to achieve gray function

## configuration

1. configuration plug-in configuration

```yaml
experimental:
  plugins:
    request-dispatch:
      moduleName: github.com/qxsugar/request-dispatch
      version: v1.0.1
```

2. route configuration

```yaml
http:
  routers:
    api:
      rule: host(`api.test.cn`)
      service: svc
      entryPoints:
        - web
      middlewares:
        - request-dispatch
  services:
    svc:
      loadBalancer:
        servers:
          - url: "https://prod.api.cn"
  middlewares:
    request-dispatch:
      plugin:
        request-dispatch:
          logLevel: DEBUG
          markHeader: X-DISPATCH
          markHosts:
            alpha:
              - https://alpha.api.cn
              - https://alpha1.api.cn
            beta:
              - https://beta.api.cn
```

## used

If the request header takes the mark header parameter, the request will be dispatch to the appropriate address

> http api.test.cn X-DISPATCH:alpha -v # the request will be dispatch to `https://alpha.api.cn` or `https://alpha1.api.cn`
>
> http api.test.cn X-DISPATCH:beta -v # the request will be dispatch to `https://beta.api.cn`
>
> http api.test.cn X-DISPATCH:whoami -v # the request will not be dispatch