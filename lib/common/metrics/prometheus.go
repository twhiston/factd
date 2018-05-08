package metrics

import "github.com/prometheus/client_golang/prometheus"

// PromErrorCount contains the number of non fatal errors that occurred during runtime
var PromErrorCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "error_total",
	Help: "Number of non fatal errors during runtime.",
})

// PromEnabledPlugins contains the number of active plugin
var PromEnabledPlugins = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "enabled_plugins",
	Help: "Number of enabled plugin.",
})
