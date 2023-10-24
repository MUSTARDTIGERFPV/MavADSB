package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

var startTime time.Time
var COMPILE_VERSION string
var COMPILE_HOSTNAME string
var COMPILE_TIMESTAMP string
var COMPILE_USER string

func getIPInfo() (IPInfo, error) {
	url := "http://ip-api.com/json"
	resp, err := http.Get(url)
	if err != nil {
		return IPInfo{}, err
	}
	defer resp.Body.Close()

	var info IPInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return IPInfo{}, err
	}

	return info, nil
}

func statusString() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return fmt.Sprintf("%s@%s version %s (%s) up %s | built at %s by %s@%s\n\n", os.Args[0], hostname,
		COMPILE_VERSION, runtime.Version(), time.Since(startTime), COMPILE_TIMESTAMP, COMPILE_USER, COMPILE_HOSTNAME)
}

func varsString() string {
	return fmt.Sprintf("%+v", flags)
}

func removeConn(slice []net.Conn, index int) []net.Conn {
	if len(slice) < 2 {
		return []net.Conn{}
	}
	return append(slice[:index], slice[index+1:]...)
}
