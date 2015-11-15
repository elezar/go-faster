package logs

import (
	"time"

	"github.com/emicklei/go-restful"
	"github.com/pakohan/go-faster/api/controller"
	"github.com/pakohan/go-faster/api/model"
)

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
	d, err := model.GetData(request.PathParameter("series_name"), time.Time{}, time.Time{})
	controller.GetHandler(response, d, err)
}

func getSeries(request *restful.Request, response *restful.Response) {
	s := model.ListSeries()
	controller.GetHandler(response, s, nil)
}

func createSeries(request *restful.Request, response *restful.Response) {
	s, err := model.CreateSeries(request.PathParameter("series_name"))
	controller.CreateHandler(response, s, err)
}

func deleteSeries(request *restful.Request, response *restful.Response) {
	err := model.DeleteSeries(request.PathParameter("series_name"))
	controller.DeleteHandler(response, err)
}
