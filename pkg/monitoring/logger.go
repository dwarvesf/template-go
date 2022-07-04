package monitoring

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// logger handles logging to std out in JSON with tag support
type logger struct {
	// Currently unable to retrieve tags saved in uber zap logger due to its design to be quick.
	// Hence keeping a local copy of tags for other purpose such as sentry error reporting
	tags      map[string]string
	zap       *zap.Logger
	debugMode bool
}

// WithTag creates new child Monitor and adds a new tag to it. Parent Monitor remains unchanged.
// Currently uber zap allows duplicated tags.
// Ref: https://github.com/uber-go/zap/blob/1cac10bfebaa55cacfa76e277137ec4167987ec6/zapcore/json_encoder.go#L80
// However, these tags will be dedup when processed by Fluentd
func (m *Monitor) WithTag(key, value string) *Monitor {
	if m == nil {
		return nil
	}

	// create and return new monitor
	return &Monitor{
		sentryClient: m.sentryClient,
		logger:       m.logger.withTag(key, value),
	}
}

// WithTags creates a new child Monitor and adds new tags to it. Parent Monitor remains unchanged.
// Currently uber zap allows duplicated tags.
// Ref: https://github.com/uber-go/zap/blob/1cac10bfebaa55cacfa76e277137ec4167987ec6/zapcore/json_encoder.go#L80
// However, these tags will be dedup when processed by Fluentd
func (m *Monitor) WithTags(tags map[string]string) *Monitor {
	if m == nil {
		return nil
	}

	// create and return new monitor
	return &Monitor{
		sentryClient: m.sentryClient,
		logger:       m.logger.withTags(tags),
	}
}

// Debugf logs the message using debug level
func (m *Monitor) Debugf(format string, args ...interface{}) {
	if m == nil {
		return
	}
	m.logger.debugf(format, args...)
}

// Infof logs the message using info level
func (m *Monitor) Infof(format string, args ...interface{}) {
	if m == nil {
		return
	}
	m.logger.infof(format, args...)
}

// Errorf logs the message using error level and reports the error to sentry
func (m *Monitor) Errorf(err error, extraMessageFormat string, args ...interface{}) {
	if m == nil {
		return
	}
	m.logger.errorf(err, extraMessageFormat, args...)
	m.ReportError(err, m.logger.tags)
}

// DebugMode returns debugMode flag
func (m *Monitor) DebugMode() bool {
	if m == nil || m.logger == nil {
		return false
	}
	return m.logger.debugMode
}

func (l *logger) debugf(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zap.Debug(fmt.Sprintf(format, args...))
}

func (l *logger) infof(format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.zap.Info(fmt.Sprintf(format, args...))
}

func (l *logger) errorf(err error, format string, args ...interface{}) {
	if l == nil {
		return
	}
	if format != "" {
		l.zap.Error(fmt.Sprintf(format+". Err: %v", append(args, err)...))
		return
	}
	l.zap.Error(fmt.Sprintf("Err: %v", err))
}

// withTag creates a new child Logger and adds a new tag to it. Parent Logger remains unchanged.
func (l *logger) withTag(key, value string) *logger {
	if l == nil {
		return nil
	}

	clone := l.clone()

	// add tags
	clone.tags[key] = value
	clone.zap = clone.zap.With(zap.String(key, value))

	// create and return new logger
	return clone
}

// withTags creates a new child Logger and adds new tags to it. Parent Logger remains unchanged.
func (l *logger) withTags(tags map[string]string) *logger {
	if l == nil {
		return nil
	}

	clone := l.clone()

	// add new tags
	for key, value := range tags {
		clone.tags[key] = value
		clone.zap = clone.zap.With(zap.String(key, value))
	}

	return clone
}

func (l *logger) clone() *logger {
	t := make(map[string]string)

	for k, v := range l.tags {
		t[k] = v
	}

	return &logger{
		tags:      t,
		zap:       l.zap,
		debugMode: l.debugMode,
	}
}

// flush flushes any buffered log entries.
func (l *logger) flush(maxWait time.Duration) error {
	if l == nil {
		return nil
	}

	errChan := make(chan error)
	go func() {
		errChan <- l.zap.Sync()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), maxWait)
	defer cancel()
	select {
	case <-errChan:
		// NOTE: We ignore any errors here because Sync is known to fail with EINVAL
		// When logging to Stdout on certain OS's.
		//
		// Uber made the same change within the core of the logger implementation.
		// See: https://github.com/uber-go/zap/issues/328
		// See: https://github.com/influxdata/influxdb/pull/20448
		return nil
	case <-ctx.Done():
		return errors.New("timed out")
	}
}
