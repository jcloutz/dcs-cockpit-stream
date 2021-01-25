package metrics

import (
	"time"

	"github.com/armon/go-metrics"
)

type MetricType string

const (
	// Samples
	MetricSampleCaptureController MetricType = "capture-controller"
	MetricSampleViewportHandler   MetricType = "viewport-handler"
	MetricPipelineExecutionTime   MetricType = "pipeline-execution"

	// Counters
	MetricFrameCounter MetricType = "frames"
)

type Service struct {
	sink *metrics.InmemSink
}

func New() *Service {
	sink := metrics.NewInmemSink(5*time.Second, 10*time.Second)

	return &Service{
		sink: sink,
	}
}

func (m *Service) MeasureTime(start time.Time, name MetricType) {
	elapsed := time.Since(start)

	m.sink.AddSample([]string{string(name)}, float32(elapsed.Nanoseconds()))
}

func (m *Service) MeasureTimeForClient(start time.Time, name MetricType, clientId string) {
	elapsed := time.Since(start)

	m.sink.AddSample([]string{string(name)}, float32(elapsed.Nanoseconds()))
	m.sink.AddSampleWithLabels([]string{string(name)}, float32(elapsed.Nanoseconds()), []metrics.Label{{"client", clientId}})
}

func (m *Service) Count(name MetricType, amount float32) {
	m.sink.IncrCounter([]string{string(name)}, amount)
}

func (m *Service) CountClient(name MetricType, amount float32, clientId string) {
	m.sink.IncrCounter([]string{string(name)}, amount)
	m.sink.IncrCounterWithLabels([]string{string(name)}, amount, []metrics.Label{{"client", clientId}})
}

func (m *Service) Data() *Snapshot {
	return &Snapshot{
		data: m.sink.Data()[0],
	}
}
