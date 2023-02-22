# pyramidel-chain-baas

国密fabric baas

## 整体功能

### 动态fabric网络

#### 组织

- [ ] **新增组织到系统通道**

## 项目目录解构

## 端口占用情况

| serve | port |
| :------------: | :--: |
| http serve | 8080 |
| jsonrpc server | 8082 |
| docker daemon | 2375 |

## nohup启用服务

export PYCBAAS_CFG_PATH=/root/pyramidel-chain-baas/configs
nohup /root/pyramidel-chain-baas/cmd/service/serve >/root/serve.log 2>&1 &