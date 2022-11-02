package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var runStateDesc = prometheus.NewDesc(
	"tfcbadge_run_state_timestamp",
	"Current state of runs.",
	[]string{"workspace_id", "workspace_name", "trigger", "status"}, nil,
)

type MetricsCollector struct {
	Store *RunCache
}

func (c *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	runs := c.Store.List()
	for _, run := range runs {
		notification := run.Notifications[0]
		ch <- prometheus.MustNewConstMetric(
			runStateDesc,
			prometheus.CounterValue,
			float64(notification.RunUpdatedAt.Unix()),
			run.WorkspaceID,
			run.WorkspaceName,
			notification.Trigger,
			notification.RunStatus,
		)
	}
}
