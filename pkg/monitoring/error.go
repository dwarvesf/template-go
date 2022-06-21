package monitoring

import (
	"github.com/getsentry/sentry-go"
)

// ReportError will report the error
func (m *Monitor) ReportError(err error, tags map[string]string) {
	if m == nil || m.sentryClient == nil {
		return
	}
	scope := sentry.NewScope()
	scope.SetTags(tags)
	m.sentryClient.CaptureException(err, nil, scope)
}
