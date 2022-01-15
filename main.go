package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	httpDurationsHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_durations_histogram_seconds",
			Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30},
		},
		[]string{"path"},
	)
	prometheus.MustRegister(httpDurationsHistogram)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//统计延时
		start := time.Now()
		defer func() {

			end := time.Now()
			httpDurationsHistogram.With(
				prometheus.Labels{"path": r.URL.Path},
			).Observe(end.Sub(start).Seconds())
		}()

		//延时0~2秒
		rand.Seed(time.Now().UTC().UnixNano())
		delay := rand.Intn(2000)
		time.Sleep(time.Millisecond * time.Duration(delay))

		lowerCaseHeader := make(http.Header)
		for k, vals := range r.Header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
			lowerCaseHeader[strings.ToLower(k)] = vals
		}

		if envVar, has := os.LookupEnv("VERSION"); has {
			w.Header().Add("ENV_VAR", envVar)
		}

		w.WriteHeader(200)
		w.Write([]byte("OK"))

		log.Printf("%s:%d", getIP(r), 200)

	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8888", nil)

}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

