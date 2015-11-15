package model

import (
	"sync"
	"time"
)

var (
	l       sync.Mutex
	entries = []Entry{}
)

type Entry struct {
	Timestamp       time.Time `json:"timestamp"`
	DistanceM       uint64    `json:"distance_m"`
	LatencyMS       uint64    `json:"latency_ms"`
	DownloadSpeedBS uint64    `json:"download_speed_bs"`
	UploadSpeedBS   uint64    `json:"uploadx_speed_bs"`
}

type DataContainer struct {
	Data [][]interface{} `json:"data"`
}

func AddEntry(e ...Entry) (err error) {
	l.Lock()
	defer l.Unlock()

	entries = append(entries, e...)
	return
}

func GetData(from, until time.Time) (dc *DataContainer, err error) {
	l.Lock()
	defer l.Unlock()

	dc = new(DataContainer)
	for _, e := range entries {
		dc.Data = append(dc.Data, []interface{}{e.Timestamp.Unix(), e.DownloadSpeedBS, e.UploadSpeedBS})
	}
	return
}
