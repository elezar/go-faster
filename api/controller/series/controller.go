package series

import (
	"github.com/emicklei/go-restful"
	"github.com/pakohan/go-faster/api/controller"
	"github.com/pakohan/go-faster/api/model"
)

func init() {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Path("/series")

	initWebResource(ws)
	controller.AddResource(ws)
}

func initWebResource(ws *restful.WebService) {
	seriesName := ws.PathParameter("series_name", "group to log into")

	from := ws.QueryParameter("from", "timestamp of the oldest data to retrieve")
	until := ws.QueryParameter("until", "timestamp of the latest data to retrieve")
	field := ws.QueryParameter("field", "fields to add")

	ws.Route(ws.POST("").To(createSeries).
		Doc("create a new series").
		Reads(model.Series{}).
		Writes(model.Series{}))

	ws.Route(ws.GET("").To(getSeries).
		Doc("get all available log series").
		Writes([]model.Series{}))

	ws.Route(ws.GET("/{series_name}").To(getSingleSeries).
		Doc("get a single series meta data").
		Param(seriesName).
		Writes(model.Series{}))

	ws.Route(ws.PUT("/{series_name}").To(updateSeries).
		Doc("update the meta data of a single series").
		Param(seriesName).
		Reads(model.Series{}).
		Writes(model.Series{}))

	ws.Route(ws.DELETE("/{series_name}").To(deleteSeries).
		Doc("delete an existing series").
		Param(seriesName))

	ws.Route(ws.PUT("/{series_name}/data").To(addData).
		Doc("add a single data point to a series").
		Param(seriesName).
		Reads(model.Entry{}).
		Writes(model.Entry{}))

	ws.Route(ws.GET("/{series_name}/data/conversions").To(getData).
		Doc("get data").
		Param(seriesName).
		Param(from).
		Param(until).
		Param(field).
		Writes(model.DataContainer{}))

}
