package model

import (
	"time"

	"log"
)

type Field struct {
	Field, Format string
}

type Fields []Field

func (self Fields) GetRow(metaData *Series, e Entry) (d []interface{}) {
	for _, f := range self {
		var v interface{}
		switch f.Field {
		case "time", "unixtime", "timestamp":
			switch f.Format {
			case "unixtime":
				v = e.Unixtime
			case "date":
				v = time.Unix(e.Unixtime, 0).Format("2006-01-02")
			case "datetime":
				v = time.Unix(e.Unixtime, 0).Format("2006-01-02T15:04:05")
			default:
				log.Printf("unkown time format '%s'", f.Format)
			}
		case "upload_speed":
			v = fmtSpeed(f.Format, e.UploadSpeedBS)
		case "download_speed":
			v = fmtSpeed(f.Format, e.DownloadSpeedBS)
		case "expected_upload_speed":
			v = fmtSpeed(f.Format, metaData.ExpectedUploadSpeed)
		case "expected_download_speed":
			v = fmtSpeed(f.Format, metaData.ExpectedDownloadSpeed)
		case "latency", "ping":
			v = e.LatencyMS
		default:
			log.Printf("unknown field '%s'", f.Field)
		}

		d = append(d, v)
	}

	return
}

func fmtSpeed(format string, speed uint64) (v interface{}) {
	fspeed := float64(speed)
	switch format {
	case "gbit_s":
		fspeed = fspeed / 1024
		fallthrough
	case "mbit_s":
		fspeed = fspeed / 1024
		fallthrough
	case "kbit_s":
		fspeed = fspeed / 1024
	}

	v = fspeed
	return
}
