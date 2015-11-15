package model

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"sync"
)

type Series struct {
	Name string       `json:"name"`
	l    sync.RWMutex `json:"-"`
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

		entries = append(entries, *e)
	}

	err = scanner.Err()
	return
}
