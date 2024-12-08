package metrics

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	track_state    = "track_state"
	refund_state   = "refund_state"
	error_reason   = "error_reason"
	feedback_state = "feedback_state"
	page_type      = "page_type"

	dataprovider = "dataprovider"
	method       = "http_method"
	status_code  = "status_code"

	lat_lng_type = "lat_lng_type"
)

var (
	sandappconfigWpGauge                        *prometheus.GaugeVec
	sandappconfigXpConfigParamsGauge            *prometheus.GaugeVec
	xpConfigParamsGauge                         *prometheus.GaugeVec
	ratingsServiceErrorCountGauge               *prometheus.GaugeVec
	trackStatusCounter                          *prometheus.CounterVec
	httpStatusCounter                           *prometheus.CounterVec
	httpStatusRetryCounter                      *prometheus.CounterVec
	launchAPILatLngCounter                      *prometheus.CounterVec
	dataAggregatorLatencyHistogram              *prometheus.HistogramVec
	hystrixClientWithRetryLatencyHistogram      *prometheus.HistogramVec
	ozoneAuthenticatorCounter                   *prometheus.CounterVec
	feedbackStatusCounter                       *prometheus.CounterVec
	refundStatusCounter                         *prometheus.CounterVec
	pageTypeStatusCounter                       *prometheus.CounterVec
	orderHistoryEnrichmentUserIdMismatchCounter *prometheus.CounterVec
)

// {{service}}-{{dataprovider}-{{http_path}}-{{http_method}}-{{status_code}}
func init() {
	// Specifically used to Track Error in Track Crouton errors
	// Part of https://docs.google.com/document/d/19rqv_EZ3RaxdLl4tUIlp12UhOhrJPfuY9AQGq7t90vY/edit#heading=h.4bh2ke1bwtew
	trackStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_launch_food_track_status_count",
		},
		[]string{track_state, error_reason})
	prometheus.Unregister(trackStatusCounter)
	prometheus.MustRegister(trackStatusCounter)

	refundStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_launch_refund_status_count",
			Help: "Counter to track frequency of different outcomes related to launch API payment population",
		},
		[]string{refund_state})
	prometheus.Unregister(refundStatusCounter)
	prometheus.MustRegister(refundStatusCounter)

	feedbackStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_launch_food_feedback_status_count",
		},
		[]string{feedback_state})
	prometheus.Unregister(feedbackStatusCounter)
	prometheus.MustRegister(feedbackStatusCounter)

	pageTypeStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_settings_api_pagetype_status_count",
		},
		[]string{page_type})
	prometheus.Unregister(pageTypeStatusCounter)
	prometheus.MustRegister(pageTypeStatusCounter)

	sandappconfigWpGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "launch_food_wp_options_gauge",
			Help:      "Count of Wp options received from sand-app-config .",
		}, []string{},
	)
	prometheus.MustRegister(sandappconfigWpGauge)

	sandappconfigXpConfigParamsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "launch_food_xp_config_params_gauge",
			Help:      "Count of Xp config params received from sand-app-config .",
		}, []string{},
	)
	prometheus.MustRegister(sandappconfigXpConfigParamsGauge)

	xpConfigParamsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "xp_config_params_gauge",
			Help:      "Count of Xp config params received from XP.",
		},
		[]string{})
	prometheus.MustRegister(xpConfigParamsGauge)

	ratingsServiceErrorCountGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "foundation_gateway_service",
			Name:      "ratings_service_error_count_gauge",
			Help:      "Count of Empty Response received form rating service.",
		},
		[]string{})
	prometheus.MustRegister(ratingsServiceErrorCountGauge)

	launchAPILatLngCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_launch_lat_lng_count",
		},
		[]string{lat_lng_type})
	prometheus.Unregister(launchAPILatLngCounter)
	prometheus.MustRegister(launchAPILatLngCounter)

	orderHistoryEnrichmentUserIdMismatchCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_order_history_enrichment_user_id_mismatch_count",
		},
		[]string{"request_user_id", "response_user_id"},
	)
	prometheus.Unregister(orderHistoryEnrichmentUserIdMismatchCounter)
	prometheus.MustRegister(orderHistoryEnrichmentUserIdMismatchCounter)

	httpStatusCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_client_http_error_status_code_count",
		},
		[]string{dataprovider, method, status_code})
	prometheus.Unregister(httpStatusCounter)
	prometheus.MustRegister(httpStatusCounter)

	httpStatusRetryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_client_http_retry_count",
		},
		[]string{dataprovider, method, status_code})
	prometheus.Unregister(httpStatusRetryCounter)
	prometheus.MustRegister(httpStatusRetryCounter)

	dataAggregatorLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "foundation_gateway_service",
			Name:      "dataprovider_aggregator_duration_seconds",
			Help:      "A histogram of request latencies.",
			Buckets:   HistBuckets,
		},
		[]string{"dataprovider"},
	)
	prometheus.Unregister(dataAggregatorLatencyHistogram)
	prometheus.MustRegister(dataAggregatorLatencyHistogram)

	hystrixClientWithRetryLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "foundation_gateway_service",
			Name:      "hystrix_client_with_retry_latency",
			Help:      "A histogram of request latencies.",
			Buckets:   HistBuckets,
		},
		[]string{"hystrixclient"},
	)
	prometheus.Unregister(hystrixClientWithRetryLatencyHistogram)
	prometheus.MustRegister(hystrixClientWithRetryLatencyHistogram)
	ozoneAuthenticatorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "foundation_gateway_ozone_authenticator_count",
		},
		[]string{"api"},
	)
	prometheus.Unregister(ozoneAuthenticatorCounter)
	prometheus.MustRegister(ozoneAuthenticatorCounter)
}

func IncrementTrackStatusCounter(state, reason string) {
	trackStatusCounter.WithLabelValues(state, reason).Inc()
}

func IncrementRefundStatusCounter(state string) {
	refundStatusCounter.WithLabelValues(state).Inc()
}

func IncrementLaunchAPILatLngCounter(latLngType string) {
	launchAPILatLngCounter.WithLabelValues(latLngType).Inc()
}

// /private package function to test the counter functionality
func getLaunchAPILatLngMetricCounter() *prometheus.CounterVec {
	return launchAPILatLngCounter
}

func getDataProviderLatencyHistogram() *prometheus.HistogramVec {
	return dataAggregatorLatencyHistogram
}

func getHystrixClientLatencyHistogram() *prometheus.HistogramVec {
	return hystrixClientWithRetryLatencyHistogram
}

func IncrementClientHTTPStatusCounter(http_provider, http_method, http_status_code string) {
	httpStatusCounter.WithLabelValues(http_provider, http_method, http_status_code).Inc()
}

func IncrementClientHTTPRetryStatusCounter(http_provider, http_method, http_status_code string) {
	httpStatusRetryCounter.WithLabelValues(http_provider, http_method, http_status_code).Inc()
}

func InstrumentDataProviderLatencyHistogram(dataprovider string, startTime time.Time) {
	timeDiff := float64(time.Since(startTime).Seconds())
	if dataprovider == "" {
		return
	}
	dataAggregatorLatencyHistogram.WithLabelValues(dataprovider).Observe(timeDiff)
}

func InstrumentHystrixClientLatencyHistogram(dataprovider string, startTime time.Time) {
	timeDiff := float64(time.Since(startTime).Seconds())
	if dataprovider == "" {
		return
	}
	hystrixClientWithRetryLatencyHistogram.WithLabelValues(dataprovider).Observe(timeDiff)
}

func IncrementOzoneAuthenticatorCounter(apiPath string) {
	ozoneAuthenticatorCounter.WithLabelValues(apiPath).Inc()
}

func RecordSandAppConfigWpKeysCountGauge(count int) {
	sandappconfigWpGauge.WithLabelValues().Set(float64(count))
}

func RecordSandAppConfigXpConfigParamsCountGauge(count int) {
	sandappconfigXpConfigParamsGauge.WithLabelValues().Set(float64(count))
}

func RecordXpConfigParamsCountGauge(count int) {
	xpConfigParamsGauge.WithLabelValues().Set(float64(count))
}

func RecordRatingsServiceErrorCountGauge(count int) {
	ratingsServiceErrorCountGauge.WithLabelValues().Set(float64(count))
}

func IncrementFeedbackStatusCounter(state string) {
	feedbackStatusCounter.WithLabelValues(state).Inc()
}

func IncrementPageTypeStatusCounter(state string) {
	pageTypeStatusCounter.WithLabelValues(state).Inc()
}

func IncrementOrderHistoryEnrichmentUserIdMismatchCounter(requestUserId string, responseUserId string) {
	orderHistoryEnrichmentUserIdMismatchCounter.WithLabelValues(requestUserId, responseUserId).Inc()
}

var (
	mu sync.Mutex
)

type CounterInterface interface {
	GetMetricWithLabelValues(lvs ...string) (prometheus.Counter, error)
}

func IsPageTypeStatusCounterIncremented(counter CounterInterface, pageType string) bool {
	mu.Lock()
	defer mu.Unlock()
	metric, err := counter.GetMetricWithLabelValues(pageType)
	if err != nil {
		return false
	}
	metricValue := metric.Desc().String()
	return metricValue != ""
}

func ResetPageTypeStatusCounters() {
	mu.Lock()
	defer mu.Unlock()
	pageTypeStatusCounter.Reset()
}

type MockPageTypeStatusCounter struct {
	mock.Mock
}

func (m *MockPageTypeStatusCounter) GetMetricWithLabelValues(lvs ...string) (prometheus.Counter, error) {
	if m == nil {
		return nil, errors.New("MockPageTypeStatusCounter is nil")
	}
	args := make([]interface{}, len(lvs))
	for i, v := range lvs {
		args[i] = v
	}
	callArgs := m.Called(args...)
	if callArgs.Get(0) != nil {
		return callArgs.Get(0).(prometheus.Counter), callArgs.Error(1)
	}
	return nil, callArgs.Error(1)
}
