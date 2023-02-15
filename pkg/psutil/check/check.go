package check

import (
	"fmt"
	"github.com/hxx258456/pyramidel-chain-baas/pkg/psutil/localcache"
	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/docker"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"time"
)

const (
	B  = 1 << 0
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

type HostInfo struct {
	UsageStat   UsageStat   `json:"usageStat"`   // 硬盘使用情况
	InfoStat    InfoStat    `json:"infoStat"`    // 主机信息
	CpuInfoStat CpuInfoStat `json:"cpuInfoStat"` // cpu信息
	MemStat     MemStat     `json:"memStat"`     // 内存使用情况
	PacketsSent uint64      `json:"packetSent"`  // 上行实时流量
	PacketsRecv uint64      `json:"packetRecv"`  // 下行实时流量
}

//UsageStat 硬盘使用信息
type UsageStat struct {
	Total       uint64 `json:"total"`       // 硬盘总量
	Free        uint64 `json:"free"`        // 未使用的
	Used        uint64 `json:"used"`        // 使用的
	UsedPercent int    `json:"usedPercent"` // 已使用百分比
}

//InfoStat 服务操作系统信息
type InfoStat struct {
	Hostname      string `json:"hostname"`      // 主机名称
	Uptime        uint64 `json:"uptime"`        // 运行时间
	BootTime      string `json:"bootTime"`      // 开机时间
	Procs         uint64 `json:"procs"`         // 进程数量
	OS            string `json:"os"`            // 操作系统类型
	KernelVersion string `json:"kernelVersion"` // 操作系统内核版本
	KernelArch    string `json:"kernelArch"`    // 操作系统架构
	DockerNum     int    `json:"dockerNum"`     // 运行容器数量
}

// CoreInfoStat cpu核心信息
type CoreInfoStat struct {
	CPU       int32   `json:"cpu"`       // 编号
	Family    string  `json:"family"`    // 代数
	Mhz       float64 `json:"mhz"`       // 主频
	CacheSize int32   `json:"cacheSize"` // 缓存大小
	Percent   float64 `json:"percent"`   // 使用率
}

// CpuInfoStat cpu信息
type CpuInfoStat struct {
	Cores  []CoreInfoStat `json:"cores"` // 核心信息
	Load1  float64        `json:"load1"`
	Load5  float64        `json:"load5"`
	Load15 float64        `json:"load15"`
}

type MemStat struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`
}

//DiskCheck 服务器硬盘使用量
func DiskCheck() (UsageStat, error) {
	var usage = UsageStat{}
	u, err := disk.Usage("/")
	if err != nil {
		return usage, err
	}

	if err := copier.Copy(&usage, u); err != nil {
		return usage, err
	}
	return usage, nil
}

//OSCheck 内核检测 操作系统信息获取
func OSCheck() (InfoStat, error) {
	var statInfo = InfoStat{}
	info, err := host.Info()
	if err != nil {
		return statInfo, err
	}

	statInfo.Uptime = info.Uptime / (60 * 60 * 24)
	statInfo.OS = fmt.Sprintf("%s %s %s %s", info.Platform, info.OS, info.PlatformFamily, info.PlatformVersion)
	statInfo.Procs = info.Procs
	statInfo.KernelVersion = info.KernelVersion
	statInfo.KernelArch = info.KernelArch
	statInfo.Hostname = info.Hostname
	statInfo.BootTime = time.Unix(int64(info.BootTime), 0).Format("2006-01-02 15:04:05")
	if statInfo.DockerNum, err = checkDocker(); err != nil {
		return statInfo, err
	}
	return statInfo, nil
}

// CPUCheck cpu使用量
func CPUCheck() (CpuInfoStat, error) {
	var cpuInfo CpuInfoStat
	cpus, err := cpu.Info()
	if err != nil {
		return cpuInfo, err
	}

	cpuInfo.Cores = []CoreInfoStat{}

	pers, err := cpu.Percent(1*time.Second, true)
	if err != nil {
		return cpuInfo, err
	}
	for i, v := range cpus {
		cpuInfo.Cores = append(cpuInfo.Cores, CoreInfoStat{
			CPU:       v.CPU,
			CacheSize: v.CacheSize / 1024,
			Mhz:       v.Mhz / 1000,
			Percent:   pers[i],
			Family:    v.Family,
		})
	}
	a, err := load.Avg()
	if err != nil {
		return cpuInfo, err
	}
	cpuInfo.Load1 = a.Load1
	cpuInfo.Load5 = a.Load5
	cpuInfo.Load15 = a.Load15
	return cpuInfo, nil
}

// RAMCheck 内存使用量
func RAMCheck() (MemStat, error) {
	memStat := MemStat{}
	u, err := mem.VirtualMemory()
	if err != nil {
		return memStat, err
	}

	if err := copier.Copy(&memStat, u); err != nil {
		return memStat, err
	}
	memStat.Used /= MB
	memStat.Total /= MB
	memStat.Available /= MB
	return memStat, nil
}

func checkDocker() (int, error) {
	ids, err := docker.GetDockerIDList()
	if err != nil {
		return 0, err
	}
	return len(ids), nil
}

func checkIOCounters(info *HostInfo) error {
	IOCounters, err := net.IOCounters(true)
	if err != nil {
		return err
	}
	for _, v := range IOCounters {
		if v.Name == "eth0" {
			info.PacketsSent = v.PacketsSent / MB
			info.PacketsRecv = v.PacketsRecv / MB
		}
	}
	return nil
}

func CheckHost() (HostInfo, error) {
	host := HostInfo{}
	var err error
	if host.MemStat, err = RAMCheck(); err != nil {
		return host, err
	}
	if host.CpuInfoStat, err = CPUCheck(); err != nil {
		return host, err
	}
	if host.InfoStat, err = OSCheck(); err != nil {
		return host, err
	}
	if host.UsageStat, err = DiskCheck(); err != nil {
		return host, err
	}
	if err = checkIOCounters(&host); err != nil {
		return host, err
	}
	localcache.Cache.Set("info", host, 5*time.Minute)
	return host, nil
}
