package api

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemInfo 系统信息结构体
type SystemInfo struct {
	Host       HostInfo       `json:"host"`
	CPU        CPUInfo        `json:"cpu"`
	Memory     MemoryInfo     `json:"memory"`
	Disk       []DiskInfo     `json:"disk"`
	Network    NetworkInfo    `json:"network"`
	Load       LoadInfo       `json:"load"`
	Runtime    RuntimeInfo    `json:"runtime"`
}

// HostInfo 主机信息
type HostInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	Uptime          uint64 `json:"uptime"`
	BootTime        uint64 `json:"boot_time"`
}

// CPUInfo CPU信息
type CPUInfo struct {
	ModelName   string  `json:"model_name"`
	Cores       int32   `json:"cores"`
	Usage       float64 `json:"usage"`
	UsagePerCore []float64 `json:"usage_per_core"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Device     string  `json:"device"`
	Mountpoint string  `json:"mountpoint"`
	Fstype     string  `json:"fstype"`
	Total      uint64  `json:"total"`
	Free       uint64  `json:"free"`
	Used       uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	Interfaces []NetworkInterface `json:"interfaces"`
	IOStats    []NetIOStats       `json:"io_stats"`
}

// NetworkInterface 网络接口
type NetworkInterface struct {
	Name        string   `json:"name"`
	HardwareAddr string  `json:"hardware_addr"`
	MTU         int      `json:"mtu"`
	Flags       []string `json:"flags"`
	Addrs       []string `json:"addrs"`
}

// NetIOStats 网络IO统计
type NetIOStats struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

// LoadInfo 系统负载
type LoadInfo struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// RuntimeInfo 运行时信息
type RuntimeInfo struct {
	GoVersion   string `json:"go_version"`
	GOOS        string `json:"goos"`
	GOARCH      string `json:"goarch"`
	NumCPU      int    `json:"num_cpu"`
	NumGoroutine int   `json:"num_goroutine"`
}

// ProcessInfo 进程信息
type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float32 `json:"mem_percent"`
	Memory     uint64  `json:"memory"`
	Cmdline    string  `json:"cmdline"`
	CreateTime int64   `json:"create_time"`
}

// GetSystemMonitor 获取系统监控信息
func GetSystemMonitor(c *gin.Context) {
	var info SystemInfo

	// 获取主机信息
	hostInfo, _ := host.Info()
	info.Host = HostInfo{
		Hostname:        hostInfo.Hostname,
		OS:              hostInfo.OS,
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		Uptime:          hostInfo.Uptime,
		BootTime:        hostInfo.BootTime,
	}

	// 获取CPU信息
	cpuInfo, _ := cpu.Info()
	if len(cpuInfo) > 0 {
		cpuPercent, _ := cpu.Percent(0, true)
		cpuPercentPerCore, _ := cpu.Percent(0, false)
		
		info.CPU = CPUInfo{
			ModelName:    cpuInfo[0].ModelName,
			Cores:        cpuInfo[0].Cores,
			Usage:        sum(cpuPercentPerCore) / float64(len(cpuPercentPerCore)),
			UsagePerCore: cpuPercent,
		}
	}

	// 获取内存信息
	memInfo, _ := mem.VirtualMemory()
	info.Memory = MemoryInfo{
		Total:       memInfo.Total,
		Available:   memInfo.Available,
		Used:        memInfo.Used,
		UsedPercent: memInfo.UsedPercent,
	}

	// 获取磁盘信息
	partitions, _ := disk.Partitions(false)
	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		info.Disk = append(info.Disk, DiskInfo{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Fstype:      partition.Fstype,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}

	// 获取网络信息
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		var addrs []string
		for _, addr := range iface.Addrs {
			addrs = append(addrs, addr.Addr)
		}
		info.Network.Interfaces = append(info.Network.Interfaces, NetworkInterface{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr,
			MTU:          iface.MTU,
			Flags:        iface.Flags,
			Addrs:        addrs,
		})
	}

	ioStats, _ := net.IOCounters(true)
	for _, stat := range ioStats {
		info.Network.IOStats = append(info.Network.IOStats, NetIOStats{
			Name:        stat.Name,
			BytesSent:   stat.BytesSent,
			BytesRecv:   stat.BytesRecv,
			PacketsSent: stat.PacketsSent,
			PacketsRecv: stat.PacketsRecv,
		})
	}

	// 获取系统负载
	loadInfo, _ := load.Avg()
	info.Load = LoadInfo{
		Load1:  loadInfo.Load1,
		Load5:  loadInfo.Load5,
		Load15: loadInfo.Load15,
	}

	// 获取运行时信息
	info.Runtime = RuntimeInfo{
		GoVersion:    runtime.Version(),
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    info,
	})
}

// GetProcessList 获取进程列表
func GetProcessList(c *gin.Context) {
	processes, _ := process.Processes()
	var processList []ProcessInfo

	for _, p := range processes {
		name, _ := p.Name()
		statusSlice, _ := p.Status()
		var status string
		if len(statusSlice) > 0 {
			status = statusSlice[0]
		}
		cpuPercent, _ := p.CPUPercent()
		memInfo, _ := p.MemoryInfo()
		memPercent, _ := p.MemoryPercent()
		cmdline, _ := p.Cmdline()
		createTime, _ := p.CreateTime()

		processList = append(processList, ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			Status:     status,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
			Memory:     memInfo.RSS,
			Cmdline:    cmdline,
			CreateTime: createTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    processList,
	})
}

// GetDiskInfo 获取磁盘信息
func GetDiskInfo(c *gin.Context) {
	partitions, _ := disk.Partitions(true)
	var diskList []DiskInfo

	for _, partition := range partitions {
		usage, _ := disk.Usage(partition.Mountpoint)
		diskList = append(diskList, DiskInfo{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Fstype:      partition.Fstype,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Success",
		"data":    diskList,
	})
}

// sum 计算切片总和
func sum(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}