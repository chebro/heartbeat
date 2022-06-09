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
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, devices)
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
	newStats.Graph = updateGraph(diff, oldStats.Graph)
	newStats.Graph[29] = 1

	devicesMap[id] = newStats

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func handleGetDevice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	deviceInfo, ok := devicesMap[id]
	if ok {
		t, _ := template.ParseFiles("device.html")
		t.Execute(w, deviceInfo)
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
	var devices []Device
	var device Device
	for k := range devicesMap {
		device.Hostname = k
		device.Timestamp = devicesMap[k].Timestamp
		devices = append(devices, device)
	}
	sort.Slice(devices, func(i, j int) bool { return devices[i].Hostname < devices[j].Hostname })
	t, _ := template.ParseFiles("edit.html")
	t.Execute(w, devices)
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
