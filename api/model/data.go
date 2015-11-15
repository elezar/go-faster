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
	MetaData *Series         `json:"meta_data"`
	Data     [][]interface{} `json:"data"`
}

func AddEntries(name string, e ...Entry) (err error) {
	var s *Series
	if !SeriesExists(name) {
		s = &Series{Name: name}
		err = CreateSeries(name, s)
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

	dc = &DataContainer{MetaData: s}
	for _, e := range entries {
		dc.Data = append(dc.Data, []interface{}{
			e.Unixtime,
			float64(e.DownloadSpeedBS) / float64(1024*1024),
			float64(e.UploadSpeedBS) / float64(1024*1024),
			float64(s.ExpectedDownloadSpeed) / float64(1024*1024),
			float64(s.ExpectedUploadSpeed) / float64(1024*1024),
		})
	}

	return
}
