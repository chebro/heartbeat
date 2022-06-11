package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/go-chi/chi"
)

type Device struct {
	Hostname string
	Graph    string
}

type Overview struct {
	Devices []Device
	Version string
}

var gitCommitHash string
var templates = template.Must(template.ParseFiles("index.html", "device.html", "edit.html"))

func handleGetHome(w http.ResponseWriter, r *http.Request) {
	var devices []Device
	for k := range devicesMap {
		diff := (time.Now().Unix() - devicesMap[k].Timestamp) / 60

		graph := updateGraph(diff, devicesMap[k].Graph)
		g := plotGraph(graph)

		devices = append(devices, Device{
			Hostname: k,
			Graph:    g,
		})
	}
	sort.Slice(devices, func(i, j int) bool { return devices[i].Hostname < devices[j].Hostname })
	renderTemplate(w, "index.html", Overview{
		devices,
		gitCommitHash,
	})
}

func handlePostDevices(w http.ResponseWriter, r *http.Request) {
	var newStats DeviceStatsModel

	err := json.NewDecoder(r.Body).Decode(&newStats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := newStats.Hostname

	oldStats, ok := devicesMap[id]
	if !ok {
		oldStats = createDevice(id)
	}

	newStats.Timestamp = time.Now().Unix()
	diff := (newStats.Timestamp - oldStats.Timestamp) / 60
	newStats.Graph = updateGraph(diff+1, oldStats.Graph)
	newStats.Graph[29] = 1

	devicesMap[id] = newStats

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func handleGetDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	deviceInfo, ok := devicesMap[id]
	if ok {
		renderTemplate(w, "device.html", deviceInfo)
	} else {
		http.Error(w, "404 page not found", 404)
	}
}

func handleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	delete(devicesMap, id)
	w.WriteHeader(http.StatusNoContent)
}

func handleGetEdit(w http.ResponseWriter, r *http.Request) {
	var devices []string
	for k := range devicesMap {
		devices = append(devices, k)
	}
	sort.Slice(devices, func(i, j int) bool { return devices[i] < devices[j] })
	renderTemplate(w, "edit.html", devices)
}
