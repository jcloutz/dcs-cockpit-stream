package cockpit_stream

import (
	"time"

	"github.com/jcloutz/cockpit_stream/metrics"
)

const CaptureContextKey = iota

type CaptureContext struct {
	StartTime time.Time
	Metric    *metrics.Service
	ClientId  string
}
