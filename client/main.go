package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

type DeviceStats struct {
	Hostname string `json:"Hostname"`
	Platform string `json:"Platform"`
	Uptime   string `json:"Uptime"`
	Release  string `json:"Release"`
	Arch     string `json:"Arch"`
}

func getStats() io.Reader {
	h, _ := host.Info()

	device := DeviceStats{
		Hostname: h.Hostname,
		Platform: h.Platform,
		Release:  h.KernelVersion,
		Arch:     h.KernelArch,
		Uptime:   fmt.Sprint(h.Uptime),
	}

	j, _ := json.Marshal(device)

	return bytes.NewReader(j)
}

func sendRequest(method string, baseurl string, endpoint string) {
	w := getStats()

	req, _ := http.NewRequest(method, baseurl+endpoint, w)
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	fmt.Println("response Status:", res.Status)
}

func main() {
	var host, port string
	flag.StringVar(&host, "h", "localhost", "host address of the server")
	flag.StringVar(&port, "p", "8080", "server port")

	flag.Parse()
	url := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(host, port),
	}

	method := "POST"
	endpoint := "/api/devices"

	sendRequest(method, url.String(), endpoint)
	newTicker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-newTicker.C:
			sendRequest(method, url.String(), endpoint)
		}
	}
}
