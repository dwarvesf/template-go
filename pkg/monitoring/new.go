package monitoring

import (
	"context"
	"log"

	"github.com/dwarvesf/go-template/pkg/config"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
)

type contextKey string

const (
	monitorCtxKey = contextKey("monitoring_monitor")
)

// New initializes and returns a new New-Relic monitoring instance
func New(appCfg config.Config, debugMode bool) (*Monitor, error) {
	var sentryClient *sentry.Client
	var err error

	logger := initLogger(appCfg, debugMode)

	if appCfg.SentryDSN != "" {
		sentryClient, err = initSentry(appCfg, logger, appCfg.SentryDSN)
		if err != nil {
			return nil, err
		}
	}

	return &Monitor{
		sentryClient: sentryClient,
		logger:       logger,
	}, nil
}

// SentryInput holds sentry info
type SentryInput struct {
	DSN string
}

// FromContext gets the Monitor from context
func FromContext(ctx context.Context) *Monitor {
	if v, ok := ctx.Value(monitorCtxKey).(*Monitor); ok && v != nil {
		return v
	}

	// If monitor is not found in context
	// Then we create new monitor instance with only a logger (debugMode = false)
	// This is because logging is crucial and should always be available
	logger := initLogger(config.GetConfig(), false)
	logger.infof("New logger created as monitor not found in ctx")

	return &Monitor{logger: logger}
}

// SetInContext sets the Monitor in context
func SetInContext(ctx context.Context, monitor *Monitor) context.Context {
	return context.WithValue(ctx, monitorCtxKey, monitor)
}

// NewContext copies the monitor from old to a new context
// Use this when you want to use a new context but copy the monitor over from the original context
func NewContext(ctx context.Context) context.Context {
	return copyMonitorToNewCtx(ctx)
}

func copyMonitorToNewCtx(ctx context.Context) context.Context {
	return SetInContext(context.Background(), FromContext(ctx))
}

func initLogger(appCfg config.Config, debugMode bool) *logger {
	log.Println("Initializing Logger")

	logger := &logger{
		tags:      map[string]string{},
		zap:       newZapLogger(debugMode),
		debugMode: debugMode,
	}

	logger = logger.withTags(map[string]string{
		"app":        appCfg.ServiceName,
		"service":    appCfg.ServiceName,
		"env":        appCfg.Env,
		"version":    appCfg.Version,
		"request_id": uuid.NewString(),
	})

	logger.infof("Logger initialized")

	return logger
}

func initSentry(appCfg config.Config, logger *logger, dsn string) (*sentry.Client, error) {
	if dsn == "" {
		logger.infof("Sentry DSN not provided. Not using Sentry Error Reporting")
		return nil, nil
	}

	client, err := sentry.NewClient(
		sentry.ClientOptions{
			Dsn:              dsn,
			AttachStacktrace: true,
			SampleRate:       1, // send all events
			// Integrations: func(integrations []sentry.Integration) []sentry.Integration { // may need to enable this if ever we go to Sentry cloud so as to not expose our source code.
			// 	var filteredIntegrations []sentry.Integration
			// 	for _, integration := range integrations {
			// 		if integration.Name() == "ContextifyFrames" {
			// 			continue
			// 		}
			// 		filteredIntegrations = append(filteredIntegrations, integration)
			// 	}
			// 	return filteredIntegrations
			// },
			ServerName:  appCfg.ServiceName,
			Release:     appCfg.Version,
			Environment: appCfg.Env,
		},
	)
	if err != nil {
		return nil, err
	}

	logger.infof("Sentry Error Reporter initialized")

	return client, nil
}
