package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	envutil "github.com/argoproj/argo-workflows/v3/util/env"
)

const (
	nameOperationDuration               = `operation_duration_seconds`
	operationDurationDefaultBucketCount = 6
)

func addOperationDurationHistogram(_ context.Context, m *Metrics) error {
	maxOperationTimeSeconds := envutil.LookupEnvDurationOr("MAX_OPERATION_TIME", 30*time.Second).Seconds()
	operationDurationMetricBucketCount := envutil.LookupEnvIntOr("OPERATION_DURATION_METRIC_BUCKET_COUNT", operationDurationDefaultBucketCount)
	if operationDurationMetricBucketCount < 1 {
		log.Errorf("Invalid OPERATION_DURATION_METRIC_BUCKET_COUNT value of %d, setting to default %d", operationDurationMetricBucketCount, operationDurationDefaultBucketCount)
		operationDurationMetricBucketCount = operationDurationDefaultBucketCount
	}
	bucketWidth := maxOperationTimeSeconds / float64(operationDurationMetricBucketCount)
	// The buckets here are only the 'defaults' and can be overridden with configmap defaults
	return m.createInstrument(float64Histogram,
		nameOperationDuration,
		"Histogram of durations of operations",
		"s",
		withDefaultBuckets(prometheus.LinearBuckets(bucketWidth, bucketWidth, operationDurationMetricBucketCount)),
		withAsBuiltIn(),
	)
}

func (m *Metrics) OperationCompleted(ctx context.Context, durationSeconds float64) {
	m.record(ctx, nameOperationDuration, durationSeconds, instAttribs{})
}