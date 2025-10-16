package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type PerfMetric struct {
	CPUUsage      string
	MemoryUsage   string
	Storage       string
	EstimatedCost float64
}

// MeasurePerformanceAndCost measures Docker image size and runtime metrics
func MeasurePerformanceAndCost(image string) (PerfMetric, error) {
	ctx := context.Background()

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to create docker client: %v", err)
	}

	// Inspect image for size
	imgInspect, _, err := cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to inspect image: %v", err)
	}
	imageSizeMB := float64(imgInspect.Size) / (1024 * 1024)

	// Create temporary container
	containerResp, err := cli.ContainerCreate(ctx, &types.ContainerConfig{
		Image: image,
		Cmd:   []string{"sleep", "5"}, // short-lived container
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to create container: %v", err)
	}
	defer cli.ContainerRemove(ctx, containerResp.ID, types.ContainerRemoveOptions{Force: true})

	// Start container
	if err := cli.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{}); err != nil {
		return PerfMetric{}, fmt.Errorf("failed to start container: %v", err)
	}

	// Fetch container stats (non-streaming)
	stats, err := cli.ContainerStatsOneShot(ctx, containerResp.ID)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to get container stats: %v", err)
	}
	defer stats.Body.Close()

	// Parse stats JSON
	var statsJSON struct {
		CPUStats struct {
			CPUUsage struct {
				TotalUsage uint64 `json:"total_usage"`
			} `json:"cpu_usage"`
		} `json:"cpu_stats"`
		MemoryStats struct {
			Usage uint64 `json:"usage"`
		} `json:"memory_stats"`
	}

	bodyBytes, err := io.ReadAll(stats.Body)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to read stats body: %v", err)
	}

	if err := json.Unmarshal(bodyBytes, &statsJSON); err != nil {
		return PerfMetric{}, fmt.Errorf("failed to parse stats JSON: %v", err)
	}

	cpuUsage := fmt.Sprintf("%d units", statsJSON.CPUStats.CPUUsage.TotalUsage)
	memUsage := fmt.Sprintf("%dMi", statsJSON.MemoryStats.Usage/(1024*1024))

	// Estimate cost (custom formula)
	estimatedCost := imageSizeMB*0.01 + float64(statsJSON.CPUStats.CPUUsage.TotalUsage)/1e9*0.05 +
		float64(statsJSON.MemoryStats.Usage)/(1024*1024*1024)*0.01

	return PerfMetric{
		CPUUsage:      cpuUsage,
		MemoryUsage:   memUsage,
		Storage:       fmt.Sprintf("%.2fMB", imageSizeMB),
		EstimatedCost: estimatedCost,
	}, nil
}
