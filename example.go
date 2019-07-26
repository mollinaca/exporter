package main

import (
	//  "fmt"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"

	"github.com/mattn/go-pipeline" // go-pipeline : https://github.com/mattn/go-pipeline
	"unsafe"
)

// get Metrics
// sample(1)
func getMetrics() float64 {
	f := float64(12345)
	return f
}

// sample(2) get metrics via bash commands and pipe
func get443Estab() float64 {
	out, err := pipeline.Output(
		[]string{"netstat", "-n"},
		[]string{"grep", "ESTABLISHED"},
		[]string{"grep", ":443"},
		[]string{"wc", "-l"},
	)
	if err != nil {
		// PIPESTATUSの中に一つでも0以外が含まれると err : exit status ${status} になるので注意
		return float64(0)
	}
    f := *(*float64)(unsafe.Pointer(&out[0]))
	return float64(f)
}

const (
	namespace = "sampleMetric"
)

type myCollector struct{}

var (
	sampleGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "443_establish_count",
		Help:      "443_establish_count help",
	})
)

func (c myCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- sampleGauge.Desc()
}

func (c myCollector) Collect(ch chan<- prometheus.Metric) {
	value := get443Estab()
	ch <- prometheus.MustNewConstMetric(
		sampleGauge.Desc(),
		prometheus.GaugeValue,
		float64(value),
	)
}

var addr = flag.String("listen-address", ":19443", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	var c myCollector
	prometheus.MustRegister(c)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
