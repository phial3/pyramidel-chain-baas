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
| docker daemon | 2376 |

## Docker tls cli

docker --tlsverify --tlscacert=ca.pem --tlscert=client.pem --tlskey=client-key.pem -H tcp://8.142.106.130:2376 version

## pro-bing centos需要修改udp socket ping范围

sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
