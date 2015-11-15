package input

import (
	"github.com/emicklei/go-restful"
	"github.com/pakohan/go-faster/api/controller"
	"github.com/pakohan/go-faster/api/model"
)

func addData(request *restful.Request, response *restful.Response) {
	d := []model.Entry{}
	err := request.ReadEntity(&d)
	if err != nil {
		controller.Error(response, err)
		return
	}

	err = model.AddEntry(d...)
	controller.CreateHandler(response, d, err)
}
