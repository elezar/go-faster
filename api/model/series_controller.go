package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var (
	basePath  = filepath.Clean(os.Getenv("LOG_PATH")) + "/"
	regexName = regexp.MustCompile(`^[a-z]+(?:_[a-z]+)*$`)

	seriesLock sync.Mutex
	series     = map[string]*Series{}
)

func init() {
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("basepath is '%s'", basePath)

	var files []os.FileInfo
	files, err = ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	seriesLock.Lock()
	defer seriesLock.Unlock()
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		if !regexName.MatchString(f.Name()) {
			log.Fatalf("dirty directory: non series name folder found: '%s'", f.Name())
		}

		log.Printf("adding series '%s'", f.Name())
		series[f.Name()] = &Series{Name: f.Name()}
	}
}

func SeriesExists(name string) (ok bool) {
	seriesLock.Lock()
	defer seriesLock.Unlock()
	_, ok = series[name]
	return
}

func ListSeries() (ss []*Series) {
	seriesLock.Lock()
	defer seriesLock.Unlock()
	for _, s := range series {
		ss = append(ss, s)
	}

	return
}

func CreateSeries(name string) (s *Series, err error) {
	if !regexName.MatchString(name) {
		err = fmt.Errorf("error: '%s' does not match '%s'", name, regexName)
		return
	}

	err = os.Mkdir(basePath+name, os.ModePerm)
	if err != nil {
		return
	}

	s = &Series{Name: name}
	seriesLock.Lock()
	defer seriesLock.Unlock()
	series[s.Name] = s
	return
}

func DeleteSeries(name string) (err error) {
	if !regexName.MatchString(name) {
		err = fmt.Errorf("error: '%s' does not match '%s'", name, regexName)
		return
	}

	err = os.RemoveAll(basePath + name)
	if err != nil {
		return
	}

	seriesLock.Lock()
	defer seriesLock.Unlock()
	delete(series, name)
	return
}

func GetSeries(name string) (s *Series) {
	seriesLock.Lock()
	defer seriesLock.Unlock()
	s = series[name]
	return
}
