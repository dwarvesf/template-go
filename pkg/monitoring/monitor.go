package monitoring

import (
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

// DefaultFlushWait represents the default wait time for flushing
const DefaultFlushWait = 10 * time.Second // 10 sec (this flush) + 10 sec for server shutdown grace = 20 sec which is reasonable considering k8s grace is 30sec

// Monitor holds the values needed for this package
type Monitor struct {
	sentryClient *sentry.Client
	// nrApp        *newrelic.Application
	logger *logger
}

// Flush will flush all the monitor data left in the queue to the monitoring service. Note: this is a blocking call
func (m *Monitor) Flush(maxWait time.Duration) {
	if m == nil {
		return
	}
	var wg sync.WaitGroup

	// Zap
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = m.logger.flush(maxWait)
	}()

	// Sentry
	wg.Add(1)
	go func() {
		defer wg.Done()
		if m.sentryClient != nil {
			m.sentryClient.Flush(maxWait)
		}
	}()

	wg.Wait() // Worst case scenario the pod will get killed before we clear this statement but that is not a major problem.
}
