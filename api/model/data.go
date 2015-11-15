package model

import (
	"fmt"
	"time"
)

type Entry struct {
	Unixtime        int64  `json:"unixtime"`
	DistanceM       uint64 `json:"distance_m"`
	LatencyMS       uint64 `json:"latency_ms"`
	DownloadSpeedBS uint64 `json:"download_speed_bs"`
	UploadSpeedBS   uint64 `json:"upload_speed_bs"`
}

type DataContainer struct {
	Data [][]interface{} `json:"data"`
}

func AddEntries(name string, e ...Entry) (err error) {
	var s *Series
	if !SeriesExists(name) {
		s, err = CreateSeries(name)
		if err != nil {
			return
		}
	} else {
		s = GetSeries(name)
	}

	if s == nil {
		err = fmt.Errorf("series '%s' does not exist", name)
		return
	}

	err = s.AddEntries(e...)
	return
}

func GetData(name string, from, until time.Time) (dc *DataContainer, err error) {
	s := GetSeries(name)
	if s == nil {
		err = fmt.Errorf("series '%s' does not exist", name)
		return
	}

	var entries []Entry
	entries, err = s.GetEntries(from, until)
	if err != nil {
		return
	}

	dc = new(DataContainer)
	for _, e := range entries {
		dc.Data = append(dc.Data, []interface{}{e.Unixtime, e.UploadSpeedBS, e.DownloadSpeedBS})
	}

	return
}
