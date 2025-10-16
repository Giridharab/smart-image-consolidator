package metrics

import (
    "context"
    "fmt"

    "github.com/moby/moby/api/types"
    "github.com/moby/moby/client"
)

type PerfMetric struct {
	CPUUsage      string
	MemoryUsage   string
	Storage       string
	EstimatedCost float64
}

func MeasurePerformanceAndCost(image string) (PerfMetric, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to create docker client: %v", err)
	}
	cli.NegotiateAPIVersion(ctx)

	// Inspect image size
	imgInspect, _, err := cli.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to inspect image: %v", err)
	}
	imageSizeMB := float64(imgInspect.Size) / (1024 * 1024)

	// Run a temporary container to measure runtime stats
	containerResp, err := cli.ContainerCreate(ctx, &types.ContainerConfig{
		Image: image,
		Cmd:   []string{"sleep", "5"}, // run short-lived container
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to create container: %v", err)
	}
	defer cli.ContainerRemove(ctx, containerResp.ID, types.ContainerRemoveOptions{Force: true})

	if err := cli.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{}); err != nil {
		return PerfMetric{}, fmt.Errorf("failed to start container: %v", err)
	}

	// Get stats (non-streaming)
	stats, err := cli.ContainerStatsOneShot(ctx, containerResp.ID)
	if err != nil {
		return PerfMetric{}, fmt.Errorf("failed to get container stats: %v", err)
	}
	defer stats.Body.Close()

	// You can parse stats.Body JSON to get CPU/Memory usage
	// For simplicity, we'll simulate minimal metrics here
	// TODO: parse stats.Body for real-time metrics

	estimatedCost := imageSizeMB*0.01 + 0.1 // example: $ per MB + base cost

	return PerfMetric{
		CPUUsage:      "Simulated: run container to measure",
		MemoryUsage:   "Simulated: run container to measure",
		Storage:       fmt.Sprintf("%.2fMB", imageSizeMB),
		EstimatedCost: estimatedCost,
	}, nil
}
