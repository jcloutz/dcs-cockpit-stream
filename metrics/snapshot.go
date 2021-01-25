package metrics

import (
	"fmt"
	"github.com/armon/go-metrics"
	"time"
)

type Snapshot struct {
	data *metrics.IntervalMetrics
}

func (ms *Snapshot) getKey(metric MetricType, client string) string {
	return fmt.Sprintf("%s;client=%s", string(metric), client)
}

func (ms *Snapshot) GetSample(metric MetricType) *metrics.AggregateSample {
	return ms.data.Samples[string(metric)].AggregateSample
}

func (ms *Snapshot) GetSampleForClient(metric MetricType, client string) *metrics.AggregateSample {
	return ms.data.Samples[ms.getKey(metric, client)].AggregateSample
}

func (ms *Snapshot) GetCount(metric MetricType) *metrics.AggregateSample {
	return ms.data.Counters[string(metric)].AggregateSample
}

func (ms *Snapshot) GetCountForClient(metric MetricType, clientId string) *metrics.AggregateSample {
	return ms.data.Counters[ms.getKey(metric, clientId)].AggregateSample
}

type SnapshotOverview struct {
	TimeFrame        float64
	CapturedFrames   int
	MaxPossibleFPS   int
	AvgScreenCapture float64
	AvgHandleTime    float64
	AvgPipelineExec  float64
}

func (ms *Snapshot) GetOverview() *SnapshotOverview {
	maxFramerate := ms.GetCount(MetricFrameCounter).Sum / (ms.GetSample(MetricSampleCaptureController).Sum / float64(time.Second))
	return &SnapshotOverview{
		TimeFrame:        ms.GetSample(MetricSampleCaptureController).Sum,
		CapturedFrames:   int(ms.GetCount(MetricFrameCounter).Sum),
		MaxPossibleFPS:   int(maxFramerate),
		AvgScreenCapture: ms.GetSample(MetricSampleCaptureController).Mean(),
		AvgHandleTime:    ms.GetSample(MetricSampleViewportHandler).Mean(),
		AvgPipelineExec:  ms.GetSample(MetricPipelineExecutionTime).Mean(),
	}
}
