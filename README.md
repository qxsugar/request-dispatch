# Traefik Request Dispatch Plugin

The Traefik Request Dispatch Plugin is a powerful tool that allows for intelligent request dispatching based on request markers, facilitating the implementation of gray-scale functionality.

## Features

- Intelligent request routing based on request markers
- Enables gray-scale deployments, controlled feature rollouts, and system upgrades
- Easy to configure and use with Traefik

## Configuration

### Plugin Configuration

To configure the plugin, add the following to your Traefik configuration file:

```yaml
experimental:
  plugins:
    request-dispatch:
      moduleName: github.com/qxsugar/request-dispatch
      version: v1.0.1
```

### Route Configuration

Configure your routes as shown in the example below:

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

## Usage

If the request header includes the mark header parameter, the request will be dispatched to the appropriate address:

- `http api.test.cn X-DISPATCH:alpha -v` The request will be dispatched to `https://alpha.api.cn` or `https://alpha1.api.cn`
- `http api.test.cn X-DISPATCH:beta -v` The request will be dispatched to `https://beta.api.cn`
- `http api.test.cn X-DISPATCH:whoami -v` The request will not be dispatched

This plugin provides a flexible and powerful way to manage your request traffic, enabling you to control how and where requests are dispatched based on their marker headers.