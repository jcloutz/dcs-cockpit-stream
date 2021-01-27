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
	MetricSampleBandwidth         MetricType = "bandwidth"

	// Counters
	MetricFrameCounter MetricType = "frames"
)

type Service struct {
	sink *metrics.InmemSink
}

func New() *Service {
	sink := metrics.NewInmemSink(2*time.Second, 4*time.Second)

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

func (m *Service) AddSample(name MetricType, val float32) {
	m.sink.AddSample([]string{string(name)}, val)
}

func (m *Service) AddSampleForClient(name MetricType, clientId string, val float32) {
	m.sink.AddSample([]string{string(name)}, val)
	m.sink.AddSampleWithLabels([]string{string(name)}, val, []metrics.Label{{"client", clientId}})
}

func (m *Service) Count(name MetricType, amount float32) {
	m.sink.IncrCounter([]string{string(name)}, amount)
}

func (m *Service) CountClient(name MetricType, clientId string, amount float32) {
	m.sink.IncrCounter([]string{string(name)}, amount)
	m.sink.IncrCounterWithLabels([]string{string(name)}, amount, []metrics.Label{{"client", clientId}})
}

func (m *Service) Data() *Snapshot {
	return &Snapshot{
		data: m.sink.Data()[0],
	}
}
