package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	knownAircraft = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "aircraft_tracked",
		Help: "Currently known aircraft by the system",
	})
	updatesSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "aircraft_updates_sent",
		Help: "Aircraft position updates sent",
	})
	requestsSent = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_sent",
		Help: "Upstream requests sent",
	})
	requestsFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_failed",
		Help: "Upstream requests sent",
	})
	connectedClients = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "connected_clients",
		Help: "SBS clients connected",
	})
)

// Command-line flags
type flagDef struct {
	upstream struct {
		api_base         string
		refresh_interval uint
	}
	location struct {
		lat    float64
		lng    float64
		radius uint
	}
	sbs struct {
		port string
	}
	http struct {
		port string
	}
}

var flags flagDef

func init() {
	startTime = time.Now()
	flag.Lookup("alsologtostderr").Value.Set("true")
	flag.StringVar(&flags.upstream.api_base, "upstream.api_base", "https://api.adsb.one/v2", "ADSB API base URL")
	flag.UintVar(&flags.upstream.refresh_interval, "upstream.refresh_interval", 5, "Interval in seconds between API calls")
	flag.Float64Var(&flags.location.lat, "location.lat", 0.0, "User latitude")
	flag.Float64Var(&flags.location.lng, "location.lng", 0.0, "User longitude")
	flag.UintVar(&flags.location.radius, "location.radius", 250, "Radius to request data for in nautical miles")
	flag.StringVar(&flags.sbs.port, "sbs.port", "30003", "SBS Serving port")
	flag.StringVar(&flags.http.port, "http.port", "3000", "HTTP Serving port")
	flag.Parse()
	prometheus.MustRegister(knownAircraft)
	prometheus.MustRegister(updatesSent)
	prometheus.MustRegister(requestsSent)
	prometheus.MustRegister(requestsFailed)
	prometheus.MustRegister(connectedClients)
	glog.Infoln(statusString())
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/statusz", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "%s", statusString()) })
	http.HandleFunc("/varz", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "%s", varsString()) })

	// Serve SBS TCP feed
	server := SBSServer{}
	go server.start(flags.sbs.port)
	glog.Infof("Listening for SBS clients on port %s\n", flags.sbs.port)

	client := ADSBOneClient{
		server: &server,
	}
	go client.start()
	glog.Infof("Fetching ADS-B feed from %s\n", flags.upstream.api_base)

	http.HandleFunc("/adsb", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(client.lastResponse)
		if err != nil {
			r.Response.StatusCode = http.StatusInternalServerError
			fmt.Fprintf(w, "%s", err.Error())
			return
		}
		fmt.Fprintf(w, "%s", b)
	})

	// Serve our HTTP API
	go func() {
		glog.Error(http.ListenAndServe(fmt.Sprintf(":%s", flags.http.port), nil))
	}()
	glog.Infof("Listening for HTTP debug info on %s\n", flags.http.port)

	// Wait for a signal
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}
