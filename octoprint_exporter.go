package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
  "github.com/asohh/octoprint"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8081", "The address to listen on for HTTP requests.")
var octoprintAPIKey = flag.String("apikey", "", "API Key")
var octoprintHost = flag.String("host", "", "Host")

var (
	bed_temp_actual = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "bed_temp_actual",
			Help: "Temp of the Printbed",
		},
	)
	bed_temp_target = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "bed_temp_target",
			Help: "target temp of the Printbed",
		},
	)
	tool_temp_actual = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tool_temp_actual",
			Help: "Temp of the Printtool",
		},
	)
	tool_temp_target = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "tool_temp_target",
			Help: "target temp of the Printtool",
		},
	)
	job_time_left = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "job_time_left",
			Help: "time left on this job",
		},
	)
	job_time = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "job_time",
			Help: "time spend on this job",
		},
	)
	job_completion = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "job_completion",
			Help: "percentage of the job done",
		},
	)
)

func renew() {
	go func() {
		for {
			bedStatus := octoprint.BedStatus()
			jobStatus := octoprint.JobStatus()
			toolStatus := octoprint.ToolStatus()
			bed_temp_actual.Set(float64(bedStatus.Bed.Actual))
			bed_temp_target.Set(float64(bedStatus.Bed.Target))
			tool_temp_actual.Set(float64(toolStatus.Tool0.Actual))
			tool_temp_target.Set(float64(toolStatus.Tool0.Target))
			job_time_left.Set(float64(jobStatus.Progress.PrintTimeLeft))
			job_time.Set(float64(jobStatus.Progress.PrintTime))
			job_completion.Set(float64(jobStatus.Progress.Completion))
			time.Sleep(2 * time.Second)
		}
	}()
}
func main() {
	flag.Parse()
	fmt.Println("Connecting to" ,*octoprintHost)
	octoprint.SetAPIKey(*octoprintAPIKey)
	octoprint.SetHost(*octoprintHost)
	renew()
	prometheus.MustRegister(bed_temp_actual)
	prometheus.MustRegister(bed_temp_target)
	prometheus.MustRegister(tool_temp_actual)
	prometheus.MustRegister(tool_temp_target)
	prometheus.MustRegister(job_time_left)
	prometheus.MustRegister(job_time)
	prometheus.MustRegister(job_completion)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// func main() {
// 	api.SetAPIKey("F68F0CA17C7B490C938217B6EB64D9A7")
// 	job.Pause("")
// 	fmt.Println(job.Status().Job.File.Name)
// 	fmt.Println(bed.GetTemp(""))
// }
