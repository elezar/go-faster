package logs

import (
	"time"

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

func getData(request *restful.Request, response *restful.Response) {
	d, err := model.GetData(time.Time{}, time.Time{})
	controller.GetHandler(response, d, err)
}

func getSeries(request *restful.Request, response *restful.Response) {

}

func createSeries(request *restful.Request, response *restful.Response) {

}
