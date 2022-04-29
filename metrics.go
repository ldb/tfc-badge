package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
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
		ts, err := time.Parse(time.RFC3339, notification.RunUpdatedAt)
		if err != nil {
			log.Printf("error parsing timestamp of workspace %s for metric collection: %v", run.WorkspaceID, err)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			runStateDesc,
			prometheus.CounterValue,
			float64(ts.Unix()),
			run.WorkspaceID,
			run.WorkspaceName,
			notification.Trigger,
			notification.RunStatus,
		)
	}
}
