# pyramidel-chain-baas

国密fabric baas

## 整体功能

### 动态fabric网络

#### 组织

- [ ] **新增组织到系统通道**

## 项目目录解构

```json
{
  "usageStat": {
    // 硬盘使用情况
    "total": 63278391296,
    // 硬盘总量,单位为b,total/2的20次方转为mb单位
    "free": 47982944256,
    // 未使用的,单位为b,total/2的20次方转为mb单位
    "used": 12487077888,
    // 使用的,单位为b,total/2的20次方转为mb单位
    "usedPercent": 20
    // 已使用百分比
  },
  "infoStat": {
    // 主机信息
    "hostname": "k8s-worker-node2",
    // 主机名称
    "uptime": 11,
    // 运行时间,单位天
    "bootTime": "2023-01-29 16:53:56",
    // 开机时间
    "procs": 131,
    // 进程数量
    "os": "centos linux rhel 7.9.2009",
    // 操作系统类型
    "kernelVersion": "4.19.12-1.el7.elrepo.x86_64",
    // 操作系统内核版本
    "kernelArch": "x86_64",
    // 操作系统架构
    "dockerNum": 2
    // 运行容器数量
  },
  "cpuInfoStat": {
    // cpu信息
    "cores": [
      // 核心
      {
        "cpu": 0,
        // 核心编号
        "family": "6",
        // 代数
        "mhz": 2.5,
        // 主频
        "cacheSize": 35,
        // 缓存大小,单位mb
        "percent": 2.0202020221024846
        // 使用率
      },
      {
        "cpu": 1,
        "family": "6",
        "mhz": 2.5,
        "cacheSize": 35,
        "percent": 1.0101010110512423
      }
    ],
    "load1": 0.03,
    // 系统缓存记录的上1分钟cpu平均使用率
    "load5": 0.01,
    // 系统缓存记录的上5分钟cpu平均使用率
    "load15": 0
    // 系统缓存记录的上15分钟cpu平均使用率
  },
  "memStat": {
    // 内存使用情况
    "total": 3801,
    // 总内存 mb
    "available": 2781,
    // 可用内存 mb
    "used": 753,
    // 已用内存 mb
    "usedPercent": 19.81722792257566
    // 内存使用率
  }
}
```

