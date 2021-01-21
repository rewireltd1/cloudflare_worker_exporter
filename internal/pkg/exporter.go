package cloudflare_worker_exporter

import (
	"fmt"
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	fetcher *Fetcher
}

func NewExporter(fetcher *Fetcher) *Exporter {
	return &Exporter{
		fetcher: fetcher,
	}
}

const namespace = "cloudflare_worker"

// Metrics section
var (
	requestsUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "requests_up"),
		"Was the last Cloudflare Workers analytics requests request was successful.",
		nil, nil,
	)

	cpuTimeUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cpu_time_up"),
		"Was the last Cloudflare Workers analytics cpu time request was successful.",
		nil, nil,
	)

	cpuTimePercentile = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cpu_time_percentile"),
		fmt.Sprintf("Cloudflare Workers CPU Time per Percentile analytics requests request was successful."),
		[]string{"worker", "status", "percentile"}, nil,
	)

	requestsReceived = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "requests_received_total"),
		"How many requests have been received (per worker script).",
		[]string{"worker", "status"}, nil,
	)

	requestsErrors = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "errors_total"),
		"How many errors have been returned (per worker script).",
		[]string{"worker", "status"}, nil,
	)

	subRequestsErrors = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "subrequests_total"),
		"How many subrequests have been initiated (per worker script).",
		[]string{"worker", "status"}, nil,
	)
)

func (exporter *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- requestsUp
	ch <- requestsReceived
	ch <- requestsErrors
	ch <- subRequestsErrors

	ch <- cpuTimeUp
	ch <- cpuTimePercentile
}

func (exporter *Exporter) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	wg.Add(1)
	go exporter.collectRequests(ch, &wg)

	wg.Add(1)
	go exporter.collectCpuTime(ch, &wg)

	wg.Wait()
}

func (exporter *Exporter) collectRequests(ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := exporter.fetcher.FetchRequestCount()

	if err != nil || len(data.Viewer.Accounts) < 1 {
		ch <- prometheus.MustNewConstMetric(
			requestsUp, prometheus.GaugeValue, 0,
		)
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		requestsUp, prometheus.GaugeValue, 1,
	)

	results := data.Viewer.Accounts[0].WorkersInvocationsAdaptive
	for _, worker := range results {
		workerName := worker.Dimensions.ScriptName
		status := worker.Dimensions.Status

		ch <- prometheus.MustNewConstMetric(
			requestsReceived, prometheus.CounterValue, float64(worker.Sum.Requests), workerName, status,
		)

		ch <- prometheus.MustNewConstMetric(
			requestsErrors, prometheus.CounterValue, float64(worker.Sum.Errors), workerName, status,
		)

		ch <- prometheus.MustNewConstMetric(
			subRequestsErrors, prometheus.CounterValue, float64(worker.Sum.Subrequests), workerName, status,
		)
	}
}

func (exporter *Exporter) collectCpuTime(ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := exporter.fetcher.FetchCpuTime()

	if err != nil || len(data.Viewer.Accounts) < 1 {
		ch <- prometheus.MustNewConstMetric(
			cpuTimeUp, prometheus.GaugeValue, 0,
		)
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		cpuTimeUp, prometheus.GaugeValue, 1,
	)

	results := data.Viewer.Accounts[0].WorkersInvocationsAdaptive
	for _, worker := range results {
		workerName := worker.Dimensions.ScriptName
		status := worker.Dimensions.Status

		ch <- prometheus.MustNewConstMetric(
			cpuTimePercentile, prometheus.GaugeValue, float64(worker.Quantiles.CPUTimeP25), workerName, status, "25",
		)

		ch <- prometheus.MustNewConstMetric(
			cpuTimePercentile, prometheus.GaugeValue, float64(worker.Quantiles.CPUTimeP50), workerName, status, "50",
		)

		ch <- prometheus.MustNewConstMetric(
			cpuTimePercentile, prometheus.GaugeValue, float64(worker.Quantiles.CPUTimeP75), workerName, status, "75",
		)

		ch <- prometheus.MustNewConstMetric(
			cpuTimePercentile, prometheus.GaugeValue, float64(worker.Quantiles.CPUTimeP90), workerName, status, "90",
		)

		ch <- prometheus.MustNewConstMetric(
			cpuTimePercentile, prometheus.GaugeValue, float64(worker.Quantiles.CPUTimeP99), workerName, status, "99",
		)

		ch <- prometheus.MustNewConstMetric(
			cpuTimePercentile, prometheus.GaugeValue, float64(worker.Quantiles.CPUTimeP999), workerName, status, "999",
		)
	}
}
