package cockpit_stream

import (
	"time"

	mtrc "github.com/armon/go-metrics"
)

type MetricType string

var (
	ScreenCapure  MetricType = "capture-screen"
	HandleCapture MetricType = "handle-screen"
	Pipeline      MetricType = "screen-pipeline"
)

type MetricsService struct {
	sink mtrc.MetricSink
}

func NewMetrics() *MetricsService {
	sink := mtrc.NewInmemSink(5*time.Second, 10*time.Second)
	return &MetricsService{
		sink: sink,
	}
}

func (m *MetricsService) MeasureTime(start time.Time, name MetricType) {
	elapsed := time.Since(start)

	m.sink.AddSample([]string{string(name)}, float32(elapsed.Microseconds()))
}

func (m *MetricsService) MeasureTimeForClient(start time.Time, name MetricType, clientId string) {
	elapsed := time.Since(start)

	m.sink.AddSampleWithLabels([]string{string(name)}, float32(elapsed.Microseconds()), []mtrc.Label{{"client", clientId}})
}
