package logs

import (
	"github.com/emicklei/go-restful"
	"github.com/pakohan/go-faster/api/controller"
	"github.com/pakohan/go-faster/api/model"
)

func init() {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Path("/logs")

	initWebResource(ws)
	controller.AddResource(ws)
}

func initWebResource(ws *restful.WebService) {
	series := ws.PathParameter("series", "group to log into")

	from := ws.QueryParameter("from", "timestamp of the oldest data to retrieve")
	until := ws.QueryParameter("until", "timestamp of the latest data to retrieve")

	ws.Route(ws.POST("/{series}").To(addData).
		Doc("add log entries").
		Param(series).
		Reads([]model.Entry{}).
		Writes([]model.Entry{}))

	ws.Route(ws.GET("/series").To(getSeries).
		Doc("get all available log series").
		Writes([]model.Series{}))

	ws.Route(ws.POST("/series").To(createSeries).
		Doc("get all available log series").
		Writes([]model.Series{}))

	ws.Route(ws.GET("/{series}/converted").To(getData).
		Doc("get data").
		Param(series).
		Param(from).
		Param(until).
		Writes(model.DataContainer{}))
}
