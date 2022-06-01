package main

type DeviceStatsModel struct {
	Hostname  string    `json:"hostname"`
	Platform  string    `json:"platform"`
	Release   string    `json:"release"`
	Arch      string    `json:"arch"`
	Uptime    string    `json:"uptime"`
	Timestamp int64     `json:"timestamp"`
	Graph     []float64 `json:"graph"`
}

//type DevicesModel struct {
//	Name string
//	Info DeviceStatsModel
//}
//var devices []DevicesModel

var deviceStats DeviceStatsModel
var devicesMap = make(map[string]DeviceStatsModel)
