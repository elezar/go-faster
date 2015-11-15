package series

import (
	"fmt"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/pakohan/go-faster/api/controller"
	"github.com/pakohan/go-faster/api/model"
)

func createSeries(request *restful.Request, response *restful.Response) {
	s := new(model.Series)
	err := request.ReadEntity(s)
	if err != nil {
		controller.Error(response, err)
		return
	}

	err = model.CreateSeries(s.Name, s)
	controller.CreateHandler(response, s, err)
}

func getSeries(request *restful.Request, response *restful.Response) {
	s := model.ListSeries()
	controller.GetHandler(response, s, nil)
}

func getSingleSeries(request *restful.Request, response *restful.Response) {
	s := model.GetSeries(request.PathParameter("series_name"))
	controller.GetHandler(response, s, nil)
}

func updateSeries(request *restful.Request, response *restful.Response) {
	s := new(model.Series)
	err := request.ReadEntity(s)
	if err != nil {
		controller.Error(response, err)
		return
	}

	err = model.UpdateSeries(request.PathParameter("series_name"), s)
	controller.CreateHandler(response, s, err)
}

func deleteSeries(request *restful.Request, response *restful.Response) {
	err := model.DeleteSeries(request.PathParameter("series_name"))
	controller.DeleteHandler(response, err)
}

func addData(request *restful.Request, response *restful.Response) {
	d := new(model.Entry)
	err := request.ReadEntity(d)
	if err != nil {
		controller.Error(response, err)
		return
	}

	err = model.AddEntries(request.PathParameter("series_name"), *d)
	controller.CreateHandler(response, d, err)

}

func getData(request *restful.Request, response *restful.Response) {
	from, err := getTimestamp(request.QueryParameter("from"))
	if err != nil {
		controller.Error(response, err)
		return
	}

	until, err := getTimestamp(request.QueryParameter("until"))
	if err != nil {
		controller.Error(response, err)
		return
	}

	rawFields := request.Request.URL.Query()["field"]
	var fields model.Fields
	for _, f := range rawFields {
		p := strings.Split(f, ":")
		switch len(p) {
		case 1:
			fields = append(fields, model.Field{Field: p[0]})
		case 2:
			fields = append(fields, model.Field{Field: p[0], Format: p[1]})
		default:
			controller.Error(response, fmt.Errorf("malformed fields format '%s'", f))
		}
	}

	d, err := model.GetData(request.PathParameter("series_name"), from, until, fields)
	controller.GetHandler(response, d, err)
}

const (
	fmtDate     = "2006-01-02"
	fmtDatetime = "2006-01-02T15:04:05"
)

func getTimestamp(ts string) (t time.Time, err error) {
	switch len(ts) {
	case 0:
	case len(fmtDate):
		t, err = time.Parse(fmtDate, ts)
	case len(fmtDatetime):
		t, err = time.Parse(fmtDatetime, ts)
	default:
		err = fmt.Errorf("invalid date format, either '%s' or '%s' is allowed")
	}

	return
}
