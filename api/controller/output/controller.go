package output

import (
	"github.com/emicklei/go-restful"
	"github.com/pakohan/go-faster/api/controller"
)

func init() {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Path("/output")

	initWebResource(ws)
	controller.AddResource(ws)
}

func initWebResource(ws *restful.WebService) {
	from := ws.QueryParameter("from", "timestamp of the oldest data to retrieve")
	until := ws.QueryParameter("until", "timestamp of the latest data to retrieve")

	ws.Route(ws.GET("/").To(getData).
		Doc("get data").
		Param(from).
		Param(until).
		Writes(DataContainer{}))
}
