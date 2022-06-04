package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/guptarohit/asciigraph"
)

type Device struct {
	Hostname  string
	Timestamp int64
	Graph     string
}

func handleGetHome(w http.ResponseWriter, r *http.Request) {
	var devices []Device
	var device Device
	for k := range devicesMap {
		var g string
		graph := devicesMap[k].Graph
		if len(graph) == 0 {
			g = "no beats yet"
		} else {
			diff := (time.Now().Unix() - devicesMap[k].Timestamp) / 60
			if diff >= 1 {
				graph = append(graph, make([]float64, diff)...)
				if len(graph) > 30 {
					graph = graph[(len(graph) - 30):]
					graph = append(graph, make([]float64, len(graph)-30)...)
				}
			}
			g = asciigraph.Plot(graph, asciigraph.Precision(0), asciigraph.Height(1))
			reg := regexp.MustCompile("[0-9]")
			g = reg.ReplaceAllString(g, "")
			g = strings.ReplaceAll(g, "  ┼", "")
			g = strings.ReplaceAll(g, "  ┤", "")
		}
		device.Hostname = k
		device.Timestamp = devicesMap[k].Timestamp
		device.Graph = g
		devices = append(devices, device)
	}
	sort.Slice(devices, func(i, j int) bool { return devices[i].Hostname < devices[j].Hostname })
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, devices)
}

func handlePostDevices(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&deviceStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var prevTimestamp int64
	currTimestamp := time.Now().Unix()
	if devicesMap[deviceStats.Hostname].Timestamp == 0 {
		prevTimestamp = currTimestamp
	} else {
		prevTimestamp = devicesMap[deviceStats.Hostname].Timestamp
	}
	diff := (currTimestamp - prevTimestamp) / 60
	if len(devicesMap[deviceStats.Hostname].Graph) == 0 {
		deviceStats.Graph = make([]float64, 30)
		deviceStats.Graph[29] = 1
	} else {
		deviceStats.Graph = devicesMap[deviceStats.Hostname].Graph
	}
	if diff >= 1 {
		deviceStats.Graph = append(deviceStats.Graph, make([]float64, diff)...)
		deviceStats.Graph = append(deviceStats.Graph, []float64{1}...)
		if len(deviceStats.Graph) > 30 {
			deviceStats.Graph = deviceStats.Graph[(len(deviceStats.Graph) - 30):]
			deviceStats.Graph = append(deviceStats.Graph, make([]float64, len(deviceStats.Graph)-30)...)
		}
	}
	deviceStats.Timestamp = currTimestamp
	devicesMap[deviceStats.Hostname] = deviceStats
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func handleGetDevices(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	deviceInfo, ok := devicesMap[id]
	if ok {
		t, _ := template.ParseFiles("device.html")
		t.Execute(w, deviceInfo)
	} else {
		http.Error(w, "404 page not found", 404)
	}
}

/*
func handleDeviceNotFound(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		_, ok := devicesMap[id]
		if !ok {
			http.Error(w, "404 page not found", 404)
			return
		}
		next.ServeHTTP(w, r)
	})
}
*/
