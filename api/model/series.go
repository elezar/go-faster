package model

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"sync"
)

type Series struct {
	l                     sync.RWMutex `json:"-"`
	Name                  string       `json:"name"`
	ExpectedDownloadSpeed uint64       `json:"expected_download_speed"`
	ExpectedUploadSpeed   uint64       `json:"expected_upload_speed"`
}

func (self *Series) load() (err error) {
	var f *os.File
	f, err = os.Open(basePath + self.Name + "/metadata.json")
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(self)
	return
}

func (self *Series) Save() (err error) {
	self.l.Lock()
	defer self.l.Unlock()

	err = os.RemoveAll(basePath + self.Name + "/metadata.json")
	if err != nil {
		return
	}

	var f *os.File
	f, err = os.Create(basePath + self.Name + "/metadata.json")
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(self)
	return
}

func (self *Series) AddEntries(entries ...Entry) (err error) {
	self.l.Lock()
	defer self.l.Unlock()

	var f *os.File
	f, err = os.OpenFile(basePath+self.Name+"/data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, e := range entries {
		err = enc.Encode(e)
		if err != nil {
			return
		}
	}

	return
}

func (self *Series) GetEntries(from, until time.Time) (entries []Entry, err error) {
	self.l.RLock()
	defer self.l.RUnlock()

	var f *os.File
	f, err = os.OpenFile(basePath+self.Name+"/data.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		e := new(Entry)
		err = json.Unmarshal(scanner.Bytes(), e)
		if err != nil {
			return
		}

		if !from.IsZero() && from.After(time.Unix(e.Unixtime, 0)) {
			// log.Printf("1: %s > %s", from, time.Unix(e.Unixtime, 0))
			continue
		}

		if !until.IsZero() && until.Before(time.Unix(e.Unixtime, 0)) {
			// log.Printf("2: %s < %s", until, time.Unix(e.Unixtime, 0))
			continue
		}

		entries = append(entries, *e)
	}

	err = scanner.Err()
	return
}
