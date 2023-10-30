package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/golang/glog"
)

type Aircraft struct {
	hex      string
	flightID string
	altitude int
	lat      float64
	lon      float64
	gs       float64
	track    float64
	vspeed   int
	squawk   string
}

type SBSServer struct {
	broker *Broker[string]
}

func clientFunc(conn net.Conn, s *SBSServer) {
	glog.V(2).Infof("Task starting for %s", conn.RemoteAddr())
	msgCh := s.broker.Subscribe()
	for {
		data := <-msgCh
		writer := bufio.NewWriter(conn)
		glog.V(3).Infof("Sending to client %s\n", conn.RemoteAddr())
		writer.WriteString(data)
		err := writer.Flush()
		if err != nil {
			glog.Warning(err)
			s.broker.Unsubscribe(msgCh)
			conn.Close()
			break
		}
	}
	glog.Warningf("Task terminating for %s", conn.RemoteAddr())
}

func (s *SBSServer) start(port string) {
	s.broker = NewBroker[string]()
	go s.broker.Start()
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer ln.Close()
	// Update our connected client metrics
	go func() {
		for {
			connectedClients.Set(float64(s.broker.CountSubcribed()))
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		glog.Infof("Client connected: %s", conn.RemoteAddr())
		go clientFunc(conn, s)
	}
}

func (s *SBSServer) Publish(r ADSBOneResponse) {
	message := ""
	for _, ac := range r.Ac {
		if ac.getAltitude() > 0 {
			// Convert to intermediate format
			a := Aircraft{
				hex:      ac.Hex,
				flightID: ac.getFlight(),
				altitude: ac.getAltitude(),
				lat:      ac.Lat,
				lon:      ac.Lon,
				squawk:   ac.Squawk,
				gs:       ac.Gs,
				track:    ac.Track,
				vspeed:   ac.BaroRate,
			}
			// Message type 1 is required to parse callsign
			message += createMessage(1, a)
			// Message type 3 has the bulk of in-flight data
			message += createMessage(3, a)
			// Message type 4 has gs and track
			message += createMessage(4, a)
			updatesSent.Inc()
		}
	}
	glog.V(1).Infof("Generated a %d byte payload", len(message))
	s.broker.Publish(message)
}

func createMessage(transmissionType int, ac Aircraft) string {
	messageType := "MSG"
	sessionID := 5
	aircraftID := 0

	message := fmt.Sprintf("%s,%d,%d,%d,%s,%s,%s,%s,%s,%s,%s,%d,%f,%f,%f,%f,%d,%s,,,0,0,0,0\n",
		messageType, transmissionType, sessionID, aircraftID, strings.ToUpper(ac.hex), ac.flightID, getDate(),
		getTime(), getDate(), getTime(), ac.flightID, ac.altitude, ac.gs, ac.track, ac.lat, ac.lon, ac.vspeed, ac.squawk)

	//fmt.Println(message)
	return message
}

func getDate() string {
	return time.Now().Format("2006/01/02")
}

func getTime() string {
	return time.Now().Format("15:04:05")
}
