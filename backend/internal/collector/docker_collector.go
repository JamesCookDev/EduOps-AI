package collector

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerCollector struct {
	cli         *client.Client
	previousStats map[string]CPUStats // Armazena stats anteriores por container ID
}

func NewDockerCollector() (*DockerCollector, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DockerCollector{
		cli:           cli,
		previousStats: make(map[string]CPUStats),
	}, nil
}

func (d *DockerCollector) Collect() ([]ContainerMetrics, error) {
	ctx := context.Background()

	containers, err := d.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}

	var results []ContainerMetrics

	for _, container := range containers {

		stats, err := d.cli.ContainerStats(ctx, container.ID, false)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(stats.Body)
		stats.Body.Close()

		var statsJSON StatsJSON
		if err := json.Unmarshal(body, &statsJSON); err != nil {
			continue
		}

		// Usa stats anteriores armazenados localmente
		previousCPU, exists := d.previousStats[container.ID]
		var cpuPercent float64
		if exists {
			cpuPercent = d.calculateCPUPercentWithPrevious(&statsJSON.CPUStats, &previousCPU)
		} else {
			cpuPercent = 0.0 // Primeira leitura, sem baseline
		}

		// Armazena stats atuais para próxima coleta
		d.previousStats[container.ID] = statsJSON.CPUStats

		memUsage := statsJSON.MemoryStats.Usage
		memLimit := statsJSON.MemoryStats.Limit

		memPercent := float64(memUsage) / float64(memLimit) * 100

		inspect, err := d.cli.ContainerInspect(ctx, container.ID)
		if err != nil {
			continue
		}

		results = append(results, ContainerMetrics{
			ContainerID:  container.ID,
			Name:         container.Names[0],
			CPUPercent:   cpuPercent,
			MemoryUsage:  memUsage,
			MemoryLimit:  memLimit,
			MemoryPercent: memPercent,
			RestartCount: inspect.RestartCount,
			Timestamp:    time.Now(),
		})
	}

	// Limpa stats de containers que não existem mais
	currentIDs := make(map[string]bool)
	for _, c := range containers {
		currentIDs[c.ID] = true
	}
	for id := range d.previousStats {
		if !currentIDs[id] {
			delete(d.previousStats, id)
		}
	}

	return results, nil
}

func (d *DockerCollector) calculateCPUPercentWithPrevious(current *CPUStats, previous *CPUStats) float64 {
	cpuDelta := float64(current.CPUUsage.TotalUsage - previous.CPUUsage.TotalUsage)
	systemDelta := float64(current.SystemUsage - previous.SystemUsage)

	numCPUs := float64(len(current.CPUUsage.PercpuUsage))
	if numCPUs == 0 {
		numCPUs = 1
	}

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * numCPUs * 100.0
	}

	return 0.0
}

func (d *DockerCollector) Close() error {
	return d.cli.Close()
}