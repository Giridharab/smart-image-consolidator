package metrics

type PerfMetric struct {
	CPUUsage      string
	MemoryUsage   string
	Storage       string
	EstimatedCost float64
}

func MeasurePerformanceAndCost(image string) PerfMetric {
	// Simulated metrics, extend to real K8s metrics if needed
	return PerfMetric{
		CPUUsage:      "120m",
		MemoryUsage:   "150Mi",
		Storage:       "50Mi",
		EstimatedCost: 0.57,
	}
}

