package input

import (
	"github.com/emicklei/go-restful"
	"github.com/pakohan/go-faster/api/controller"
	"github.com/pakohan/go-faster/api/model"
)

func init() {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Path("/input")

	initWebResource(ws)
	controller.AddResource(ws)
}

func initWebResource(ws *restful.WebService) {
	ws.Route(ws.POST("/").To(addData).
		Doc("add a log entry").
		Reads([]model.Entry{}).
		Writes([]model.Entry{}))
}
