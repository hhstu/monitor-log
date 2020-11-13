package main

/*
    参考自：
	https://prometheus.io/docs/concepts/metric_types/

*/

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math"
	"net/http"
	"time"
)

func init() {
	// Metrics have to be registered to be exposed:
	//prometheus.MustRegister(uptime)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(hdUsage)
	prometheus.MustRegister(hdFailures)
	prometheus.MustRegister(apiAll)
	prometheus.MustRegister(apiEach)
}
func main() {
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":21112", nil)
}

var (
	// 简单扶摇直上的累加值
	uptime = promauto.NewCounter(prometheus.CounterOpts{
		Name: "lc_demo_app_uptime",
		Help: "服务器运行时间",
	})

	// 简单扶摇直上的累加值模拟，带个label，例如：主机每个磁盘的异常次数
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lc_hd_errors_total",
			Help: "磁盘错误次数",
		},
		[]string{"device"},
	)

	//简单的可上可下的变量值，最常用的吧。例如：主机整体cpu利用率
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "lc_cpu_usage",
		Help: "cpu使用率",
	})
	//简单的可上可下的变量值，带个label 例如：主机每个磁盘的利用率
	hdUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "lc_hd_usage",
		Help: "磁盘利用率",
	}, []string{"device"})

	// 直方图，例如请求响应时间 p99 p90 p95..
	apiAll = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "lc_api_request_time_all",
		Help:    "接口延迟all",                          // Sorry, we can't measure how badly it smells.
		Buckets: prometheus.LinearBuckets(20, 5, 5), // 5 buckets, each 5 centigrade wide.
	})

	// 直方图带个label，例如请求响应时间 p99 p90 p95..单我想细粒度到每个api
	apiEach = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "lc_api_request_time_each",
		Help:    "接口延迟each",                         // .
		Buckets: prometheus.LinearBuckets(20, 5, 5), // 5 buckets, each 5 centigrade wide.
	}, []string{"api"})
)

func recordMetrics() {
	// 累加值模拟
	go func() {
		for {
			// 累加值
			uptime.Inc()
			// 带label的累加值
			hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
			hdFailures.With(prometheus.Labels{"device": "/dev/sdb"}).Inc()
			hdFailures.With(prometheus.Labels{"device": "/dev/sdb"}).Inc()
			hdFailures.With(prometheus.Labels{"device": "/dev/sdc"}).Inc()
			hdFailures.With(prometheus.Labels{"device": "/dev/sdc"}).Inc()
			hdFailures.With(prometheus.Labels{"device": "/dev/sdc"}).Inc()
			hdFailures.With(prometheus.Labels{"device": "/dev/sdc"}).Inc()
			time.Sleep(1 * time.Second)
		}
	}()
	// 变量值
	go func() {
		for {
			for i := 0; i < 100; i++ {
				tmp := float64(i)
				cpuUsage.Set(tmp)
				time.Sleep(1 * time.Second)
			}
		}
	}()
	// 变量值带label
	go func() {
		for {
			for i := 9; i < 100; i++ {
				tmp := float64(i / 100)
				hdUsage.With(prometheus.Labels{"device": "/dev/sda"}).Set(tmp)
				time.Sleep(1 * time.Second)
			}
			for i := 10; i < 100; i++ {
				tmp := float64(i / 100)
				hdUsage.With(prometheus.Labels{"device": "/dev/sdb"}).Set(tmp)
				time.Sleep(1 * time.Second)
			}
			for i := 5; i < 100; i++ {
				tmp := float64(i / 100)
				hdUsage.With(prometheus.Labels{"device": "/dev/sdc"}).Set(tmp)
				time.Sleep(1 * time.Second)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {

			for i := 0; i < 1000; i++ {
				apiAll.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
			}
			time.Sleep(3 * time.Second)
		}
	}()

}
