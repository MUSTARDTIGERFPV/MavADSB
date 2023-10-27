package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
)

type ADSBOneClient struct {
	server       *SBSServer
	lastResponse ADSBOneResponse
}

func (c *ADSBOneClient) start() {
	// Get the lat/lng of the user from flags, or their IP geolocation if not set
	lat := flags.location.lat
	lng := flags.location.lng
	if flags.location.lat == 0.0 && flags.location.lng == 0.0 {
		ipInfo, err := getIPInfo()
		if err != nil {
			glog.Error(err)
		} else {
			lat = ipInfo.Lat
			lng = ipInfo.Lon
		}
	}

	// First iteration of the timer should fire immediately
	timer := time.NewTimer(time.Duration(0))

	for {
		<-timer.C
		resp, err := getNearby(lat, lng, flags.location.radius)
		requestsSent.Inc()
		if err != nil {
			glog.Errorf("Failed to fetch upstream: %+v\n", err)
			requestsFailed.Inc()
			// Delay more next time if we timed out
			timer.Reset(time.Duration(flags.upstream.refresh_interval*2) * time.Second)
		} else {
			glog.Infof("Completed fetch from the upstream API. Got %d aircraft.\n", len(resp.Ac))
			knownAircraft.Set(float64(len(resp.Ac)))
			c.lastResponse = resp
			go c.server.sendData(resp)
			// Run our next iteration delayed by refresh_interval seconds
			timer.Reset(time.Duration(flags.upstream.refresh_interval) * time.Second)
		}
	}
}

func getNearby(lat float64, lng float64, radius uint) (ADSBOneResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/point/%f/%f/%d", flags.upstream.api_base, lat, lng, radius))
	if err != nil {
		return ADSBOneResponse{}, err
	}
	var info ADSBOneResponse
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return info, err
	}
	return info, nil
}
