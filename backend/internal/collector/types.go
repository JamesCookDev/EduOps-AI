package collector

import "time"

type ContainerMetrics struct {
	ContainerID  string
	Name         string
	CPUPercent   float64
	MemoryUsage  uint64
	MemoryLimit  uint64
	MemoryPercent float64
	RestartCount int
	Timestamp    time.Time
}

// Estruturas para estatísticas do Docker
type CPUUsage struct {
	TotalUsage        uint64   `json:"total_usage"`
	PercpuUsage       []uint64 `json:"percpu_usage,omitempty"`
	UsageInKernelmode uint64   `json:"usage_in_kernelmode"`
	UsageInUsermode   uint64   `json:"usage_in_usermode"`
}

type CPUStats struct {
	CPUUsage    CPUUsage `json:"cpu_usage"`
	SystemUsage uint64   `json:"system_cpu_usage,omitempty"`
	OnlineCPUs  uint32   `json:"online_cpus,omitempty"`
}

type MemoryStats struct {
	Usage uint64 `json:"usage"`
	Limit uint64 `json:"limit"`
}

type StatsJSON struct {
	CPUStats    CPUStats    `json:"cpu_stats,omitempty"`
	PreCPUStats CPUStats    `json:"precpu_stats,omitempty"`
	MemoryStats MemoryStats `json:"memory_stats,omitempty"`
}