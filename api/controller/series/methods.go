package series

import (
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
	d, err := model.GetData(request.PathParameter("series_name"), time.Time{}, time.Time{})
	controller.GetHandler(response, d, err)
}
