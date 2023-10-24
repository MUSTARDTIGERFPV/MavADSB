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
	squawk   string
}

type SBSServer struct {
	clients []net.Conn
}

func (s *SBSServer) start(port string) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		s.clients = append(s.clients, conn)
		glog.Infof("Client connected: %s", conn.RemoteAddr())
	}
}

func (s *SBSServer) sendData(data ADSBOneResponse) {
	connectedClients.Set(float64(len(s.clients)))
	for i := 0; i < len(s.clients); i++ {
		conn := s.clients[i]
		writer := bufio.NewWriter(conn)

		for _, ac := range data.Ac {
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
				}
				message := ""
				// Message type 1 is required to parse callsign
				message = createMessage(1, a)
				writer.WriteString(message)
				// Message type 3 has the bulk of in-flight data
				message = createMessage(3, a)
				writer.WriteString(message)
				// Message type 4 has gs and track
				message = createMessage(4, a)
				writer.WriteString(message)
				err := writer.Flush()
				if err != nil {
					glog.Warning(err)
					glog.Warningf("Removing client %s\n", conn.RemoteAddr())
					s.clients = removeConn(s.clients, i)
					break
				}
				updatesSent.Inc()
			}
		}

	}
}

func createMessage(transmissionType int, ac Aircraft) string {
	messageType := "MSG"
	sessionID := 5
	aircraftID := 0

	message := fmt.Sprintf("%s,%d,%d,%d,%s,%s,%s,%s,%s,%s,%s,%d,%f,%f,%f,%f,,%s,,,0,0,0,0\n",
		messageType, transmissionType, sessionID, aircraftID, strings.ToUpper(ac.hex), ac.flightID, getDate(),
		getTime(), getDate(), getTime(), ac.flightID, ac.altitude, ac.gs, ac.track, ac.lat, ac.lon, ac.squawk)

	//fmt.Println(message)
	return message
}

func getDate() string {
	return time.Now().Format("2006/01/02")
}

func getTime() string {
	return time.Now().Format("15:04:05")
}
